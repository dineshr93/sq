// SPDX-FileCopyrightText: 2026 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	spdx3 "github.com/spdx/tools-golang/spdx/v3/v3_0"
)

// spdx3ContextPrefix is the JSON-LD context namespace used by SPDX 3.x
// documents. A JSON file whose top-level @context starts with this prefix is
// an SPDX 3.x document and must not be parsed by the 2.x loader.
const spdx3ContextPrefix = "https://spdx.org/rdf/"

// detectSPDXVersion inspects the top-level JSON-LD @context of a JSON SBOM.
//
// It returns ("", nil) for documents without an SPDX 3.x style @context
// (handled by the SPDX 2.x loader, including non-JSON input), ("3.0.1", nil)
// for SPDX 3.x documents of the version sq supports, and an error for SPDX
// 3.x documents of any other version.
func detectSPDXVersion(data []byte) (string, error) {
	var header struct {
		Context json.RawMessage `json:"@context"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return "", nil // not valid JSON; let the 2.x loader report it
	}
	if len(header.Context) == 0 {
		return "", nil
	}
	// @context is normally a string; some documents use an array of strings.
	var contexts []string
	if err := json.Unmarshal(header.Context, &contexts); err != nil {
		var single string
		if err := json.Unmarshal(header.Context, &single); err != nil {
			return "", nil // unrecognised @context shape
		}
		contexts = []string{single}
	}
	var spdxContext string
	for _, c := range contexts {
		if strings.HasPrefix(c, spdx3ContextPrefix) {
			spdxContext = c
			break
		}
	}
	if spdxContext == "" {
		return "", nil // other JSON-LD document, not SPDX
	}
	version := strings.TrimPrefix(spdxContext, spdx3ContextPrefix)
	if i := strings.Index(version, "/"); i >= 0 {
		version = version[:i]
	}
	if version != spdx3.Version {
		return "", fmt.Errorf("SPDX JSON-LD context version %q is not supported (sq supports SPDX %s only)", version, spdx3.Version)
	}
	return version, nil
}

// loadV3 populates the SPDX struct from an SPDX 3.0.1 JSON-LD document,
// mapping the subset of 3.0 that has a 2.x equivalent:
//
//   - licence-bearing relationships (hasConcludedLicense, hasDeclaredLicense)
//     are folded into the 2.x licence columns of the from-element;
//   - every other relationship becomes a 2.x relationship row, with the 3.0
//     type IRI converted to SCREAMING_SNAKE_CASE (dependsOn -> DEPENDS_ON);
//     lifecycle-scoped relationships are flattened (scope dropped);
//   - 3.0 `generates` becomes the 2.x GENERATED_FROM with swapped endpoints
//     (3.0: "B generates BOM"; 2.x: "BOM GENERATED_FROM B");
//   - elements and profiles with no 2.x representation are counted in
//     SkippedElements / collected in Profiles for display.
func (s *SPDX) loadV3(data []byte) error {
	doc, injected, perr := parseV3(data)
	if perr != nil {
		return fmt.Errorf("parsing SPDX %s document: %w", spdx3.Version, perr)
	}
	// Injected stubs lack creationInfo by construction, so rc4 validation
	// would warn about synthesised elements and drown real findings.
	if len(injected) == 0 {
		if err := doc.Validate(false); err != nil {
			fmt.Fprintf(os.Stderr, "warning: document failed SPDX %s validation: %v\n", spdx3.Version, err)
		}
	}

	s.Spdxid = doc.ID
	s.Name = doc.Name
	s.SpdxVersion = "SPDX-" + spdx3.Version
	s.DataLicense = licenseString(doc.DataLicense)
	// 3.0 documents are identified by their URI; there is no separate namespace.
	s.DocumentNamespace = doc.ID

	if ci, ok := doc.CreationInfo.(*spdx3.CreationInfo); ok {
		s.CreationInfo.Created = ci.Created
		creators := make([]string, 0, len(ci.CreatedBy)+len(ci.CreatedUsing))
		for _, a := range ci.CreatedBy {
			creators = append(creators, agentString(a))
		}
		for _, t := range ci.CreatedUsing {
			creators = append(creators, toolString(t))
		}
		s.CreationInfo.Creators = creators
	}

	for _, pc := range doc.ProfileConformances {
		if name := iriTail(pc.GetID()); name != "" {
			s.Profiles = append(s.Profiles, name)
		}
	}

	// Pass 1: relationships. Rows are emitted directly; licence relationships
	// additionally record values for the element pass below.
	concluded := map[string]string{}
	declared := map[string][]string{}
	for _, e := range doc.Elements {
		var typeIRI, from string
		var tos []spdx3.AnyElement
		switch r := e.(type) {
		case *spdx3.Relationship:
			typeIRI, from = r.Type.GetID(), elementID(r.From)
			tos = r.To
		case *spdx3.LifecycleScopedRelationship:
			typeIRI, from = r.Type.GetID(), elementID(r.From)
			tos = r.To
		default:
			continue
		}
		if from == "" {
			continue
		}
		name := relationshipName(typeIRI)
		if iriTail(typeIRI) == "generates" {
			// 2.x orientation: generated artifact GENERATED_FROM origin.
			name = "GENERATED_FROM"
			for _, to := range tos {
				if id := elementID(to); id != "" {
					s.Relationships = append(s.Relationships, Relationships{
						SpdxElementID:      id,
						RelationshipType:   name,
						RelatedSpdxElement: from,
					})
				}
			}
			continue
		}
		for _, to := range tos {
			id := elementID(to)
			if id == "" {
				continue
			}
			switch name {
			case "HAS_CONCLUDED_LICENSE":
				if li, ok := to.(spdx3.AnyLicenseInfo); ok {
					if ls := licenseString(li); ls != "" {
						concluded[from] = ls
					}
				}
			case "HAS_DECLARED_LICENSE":
				if li, ok := to.(spdx3.AnyLicenseInfo); ok {
					if ls := licenseString(li); ls != "" {
						declared[from] = append(declared[from], ls)
					}
				}
			case "DESCRIBES":
				s.DocumentDescribes = append(s.DocumentDescribes, id)
			}
			s.Relationships = append(s.Relationships, Relationships{
				SpdxElementID:      from,
				RelationshipType:   name,
				RelatedSpdxElement: id,
			})
		}
	}

	// Pass 2: elements.
	s.SkippedElements = map[string]int{}
	for _, e := range doc.Elements {
		if sa, ok := e.(*spdx3.SoftwareAgent); ok {
			if _, isStub := injected[sa.ID]; isStub {
				// Placeholder synthesised for a cross-document reference.
				s.SkippedElements["ExternalElement"]++
				continue
			}
		}
		switch v := e.(type) {
		case *spdx3.Package:
			p := Packages{
				Spdxid:           v.ID,
				Name:             v.Name,
				VersionInfo:      v.Version,
				DownloadLocation: string(v.DownloadLocation),
				Homepage:         string(v.HomePage),
				CopyrightText:    v.CopyrightText,
				SourceInfo:       v.SourceInfo,
				Description:      v.Description,
				Checksums:        checksums(v.VerifiedUsing),
			}
			// rc4 unifies Core/PackageVerificationCode (a 2.x verification
			// code, not a checksum) into the same verifiedUsing list.
			for _, im := range v.VerifiedUsing {
				if vc, ok := im.(*spdx3.PackageVerificationCode); ok {
					p.PackageVerificationCode.PackageVerificationCodeValue = vc.HashValue
				}
			}
			if ls, ok := concluded[v.ID]; ok {
				p.LicenseConcluded = ls
			}
			if ds, ok := declared[v.ID]; ok && len(ds) > 0 {
				p.LicenseDeclared = strings.Join(ds, " AND ")
			}
			if v.SuppliedBy != nil {
				p.Supplier = agentString(v.SuppliedBy)
			}
			for _, a := range v.OriginatedBy {
				p.Originator = agentString(a)
				break
			}
			if purpose := v.PrimaryPurpose.GetID(); purpose != "" {
				p.PrimaryPackagePurpose = strings.ToUpper(iriTail(purpose))
			}
			if purl := string(v.PackageURL); purl != "" {
				p.ExternalRefs = append(p.ExternalRefs, ExternalRefs{
					ReferenceCategory: "PACKAGE-MANAGER",
					ReferenceType:     "purl",
					ReferenceLocator:  purl,
				})
			}
			s.Packages = append(s.Packages, p)
		case *spdx3.File:
			f := Files{
				Spdxid:        v.ID,
				FileName:      v.Name,
				CopyrightText: v.CopyrightText,
				Checksums:     checksums(v.VerifiedUsing),
			}
			if ls, ok := concluded[v.ID]; ok {
				f.LicenseConcluded = ls
			}
			if ds, ok := declared[v.ID]; ok {
				f.LicenseInfoInFiles = append(f.LicenseInfoInFiles, ds...)
			}
			s.Files = append(s.Files, f)
		case *spdx3.CustomLicense:
			s.HasExtractedLicensingInfos = append(s.HasExtractedLicensingInfos, HasExtractedLicensingInfos{
				LicenseID:     licenseID(v.ID),
				ExtractedText: v.Text,
				Name:          v.Name,
			})
		case *spdx3.Document,
			*spdx3.Person, *spdx3.Organization, *spdx3.SoftwareAgent, *spdx3.Agent,
			*spdx3.Tool,
			*spdx3.ListedLicense, *spdx3.LicenseExpression,
			*spdx3.ConjunctiveLicenseSet, *spdx3.DisjunctiveLicenseSet,
			*spdx3.Relationship, *spdx3.LifecycleScopedRelationship:
			// Consumed elsewhere (creators, licence columns, extracted texts)
			// or the document itself.
		default:
			s.SkippedElements[typeName(e)]++
		}
	}

	return nil
}

// checksums converts 3.0 verifiedUsing hash methods into 2.x checksums.
func checksums(methods spdx3.IntegrityMethodList) []Checksums {
	var out []Checksums
	for _, im := range methods {
		if h, ok := im.(*spdx3.Hash); ok {
			out = append(out, Checksums{
				Algorithm:     strings.ToUpper(iriTail(h.Algorithm.GetID())),
				ChecksumValue: h.Value,
			})
		}
	}
	return out
}

// licenseString renders a 3.0 licence expression tree as a 2.x licence
// expression string. Unknown licence kinds keep their IRI verbatim so no
// licence data is silently dropped.
func licenseString(l spdx3.AnyLicenseInfo) string {
	switch v := l.(type) {
	case nil:
		return ""
	case *spdx3.LicenseExpression:
		return v.LicenseExpression
	case *spdx3.ListedLicense:
		return licenseID(v.ID)
	case *spdx3.CustomLicense:
		return licenseID(v.ID)
	case *spdx3.ConjunctiveLicenseSet:
		return joinLicenses(v.Members, " AND ")
	case *spdx3.DisjunctiveLicenseSet:
		return joinLicenses(v.Members, " OR ")
	default:
		id := v.GetID()
		if tail := iriTail(id); tail == "NOASSERTION" {
			return "NOASSERTION"
		}
		return id
	}
}

func joinLicenses(members spdx3.LicenseInfoList, sep string) string {
	parts := make([]string, 0, len(members))
	for _, m := range members {
		if ls := licenseString(m); ls != "" {
			parts = append(parts, ls)
		}
	}
	return strings.Join(parts, sep)
}

// agentString renders a 3.0 agent as a 2.x creator string.
func agentString(a spdx3.AnyAgent) string {
	switch v := a.(type) {
	case *spdx3.Person:
		return "Person: " + orID(v.Name, v.ID)
	case *spdx3.Organization:
		return "Organization: " + orID(v.Name, v.ID)
	case *spdx3.SoftwareAgent:
		return "Tool: " + orID(v.Name, v.ID)
	default:
		return "Tool: " + iriTail(a.GetID())
	}
}

// toolString renders a 3.0 creation tool as a 2.x creator string.
func toolString(t spdx3.AnyTool) string {
	if tool, ok := t.(*spdx3.Tool); ok && tool.Name != "" {
		return "Tool: " + tool.Name
	}
	return "Tool: " + iriTail(t.GetID())
}

// relationshipName converts a 3.0 relationship type IRI to its 2.x-style
// SCREAMING_SNAKE_CASE name (dependsOn -> DEPENDS_ON, hasDynamicLink ->
// HAS_DYNAMIC_LINK). Already-uppercase names and unknown types keep the IRI
// tail verbatim.
func relationshipName(typeIRI string) string {
	tail := iriTail(typeIRI)
	if tail == "" {
		return ""
	}
	sawLower := false
	for _, r := range tail {
		if unicode.IsLower(r) {
			sawLower = true
			break
		}
	}
	if !sawLower {
		return tail
	}
	var b strings.Builder
	for i, r := range tail {
		if unicode.IsUpper(r) && i > 0 {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}

// elementID returns the ID of a possibly-nil element reference.
func elementID(e spdx3.AnyElement) string {
	if e == nil {
		return ""
	}
	return e.GetID()
}

// iriTail returns the fragment (after '#') or last path segment of an IRI.
func iriTail(iri string) string {
	if i := strings.LastIndexByte(iri, '#'); i >= 0 && i+1 < len(iri) {
		return iri[i+1:]
	}
	if i := strings.LastIndexByte(iri, '/'); i >= 0 {
		return iri[i+1:]
	}
	return iri
}

// licenseID extracts a 2.x-style licence id (LicenseRef-... or listed id)
// from a 3.0 licence element IRI.
func licenseID(iri string) string {
	return iriTail(iri)
}

func orID(name, id string) string {
	if name != "" {
		return name
	}
	return iriTail(id)
}

// typeName returns a short human-readable name for a 3.0 element's concrete
// Go type (e.g. "*v3_0.AIPackage" -> "AIPackage").
func typeName(e any) string {
	t := fmt.Sprintf("%T", e)
	if i := strings.LastIndexByte(t, '.'); i >= 0 {
		return t[i+1:]
	}
	return t
}

// danglingIDRe extracts the unresolvable IDs rc4 lists in its
// "unable to find viable external IRI for" errors.
var danglingIDRe = regexp.MustCompile(`for ID: (\S+)`)

// parseV3 parses an SPDX 3.0.1 document with tools-golang rc4.
//
// rc4 has no placeholder type for relationship endpoints defined in other
// documents (a valid SPDX 3 construct), and aborts naming each dangling ID.
// When that happens, stub definitions for exactly those IDs are injected into
// @graph as SoftwareAgent placeholders and the document is re-parsed, so
// cross-document relationships survive as rows instead of killing the load.
// The returned set holds the injected IDs so loadV3 can report them as
// ExternalElement instead of losing them.
func parseV3(data []byte) (*spdx3.Document, map[string]bool, error) {
	doc := &spdx3.Document{}
	err := doc.FromJSON(bytes.NewReader(data))
	if err == nil {
		return doc, nil, nil
	}
	ids := map[string]bool{}
	for _, m := range danglingIDRe.FindAllStringSubmatch(err.Error(), -1) {
		ids[m[1]] = true
	}
	if len(ids) == 0 {
		return nil, nil, err
	}
	data2, ierr := injectPlaceholders(data, ids)
	if ierr != nil {
		return nil, nil, err
	}
	doc2 := &spdx3.Document{}
	if err2 := doc2.FromJSON(bytes.NewReader(data2)); err2 != nil {
		return nil, nil, fmt.Errorf("%w (re-parse with %d injected external placeholders failed: %v)", err, len(ids), err2)
	}
	return doc2, ids, nil
}

// injectPlaceholders appends SoftwareAgent stub definitions for ids to the
// document's @graph. The @id keyword (not the compact "id" alias, whose
// expansion the offline context does not guarantee) is what makes rc4 treat
// the relationship endpoints as internal references. The graph came from an
// unmarshal, so re-marshalling it cannot fail.
func injectPlaceholders(data []byte, ids map[string]bool) ([]byte, error) {
	var docMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &docMap); err != nil {
		return nil, err
	}
	raw, ok := docMap["@graph"]
	if !ok {
		return nil, fmt.Errorf("no @graph to extend")
	}
	var graph []any
	if err := json.Unmarshal(raw, &graph); err != nil {
		return nil, err
	}
	for id := range ids {
		graph = append(graph, map[string]string{"@id": id, "type": "SoftwareAgent"})
	}
	docMap["@graph"], _ = json.Marshal(graph)
	return json.Marshal(docMap)
}
