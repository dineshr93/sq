// SPDX-FileCopyrightText: 2026 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Fixtures vendored from github.com/spdx/spdx-examples (master):
//   example3-bin.spdx3.json          software/example3  (binary + BOM, lifecycle rels)
//   examplemaven-enriched.spdx3.json software/example14 (maven, purl, orgs)
//   simplehtr-example.spdx3.json     ai/example01       (AI/dataset profile elements)
//   dataset-example01.spdx3.json     dataset/example01  (dataset profile)
// and spdx/tools-golang v0.6.0-rc4 spdx/v3/v3_0/testdata/test.json:
//   rc4-compact-dialect.json         (compact JSON-LD dialect)

func loadFile(t *testing.T, name string) *SPDX {
	t.Helper()
	s := new(SPDX)
	if err := s.Load(filepath.Join("testdata", name)); err != nil {
		t.Fatalf("Load(%s): %v", name, err)
	}
	return s
}

func TestDetectSPDXVersion(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{"2.x document", `{"SPDXID":"SPDXRef-DOCUMENT","spdxVersion":"SPDX-2.3"}`, "", false},
		{"no context", `{"name":"x"}`, "", false},
		{"other jsonld", `{"@context":"https://schema.org"}`, "", false},
		{"3.0.1", `{"@context":"https://spdx.org/rdf/3.0.1/spdx-context.jsonld"}`, "3.0.1", false},
		{"3.0.1 in array", `{"@context":["https://schema.org","https://spdx.org/rdf/3.0.1/spdx-context.jsonld"]}`, "3.0.1", false},
		{"3.0.0", `{"@context":"https://spdx.org/rdf/3.0.0/spdx-context.jsonld"}`, "", true},
		{"4.0.0", `{"@context":"https://spdx.org/rdf/4.0.0/spdx-context.jsonld"}`, "", true},
		{"invalid json", `not json`, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectSPDXVersion([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Fatalf("detectSPDXVersion() err = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("detectSPDXVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestLoadV3UnsupportedVersion asserts a hard error for a 3.x context whose
// version segment is not 3.0.1.
func TestLoadV3UnsupportedVersion(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "future.spdx3.json")
	body := `{"@context":"https://spdx.org/rdf/9.9.9/spdx-context.jsonld","id":"urn:x","name":"future"}`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	s := new(SPDX)
	err := s.Load(p)
	if err == nil {
		t.Fatal("expected error for unsupported SPDX 3.x version, got nil")
	}
	if !strings.Contains(err.Error(), "9.9.9") {
		t.Fatalf("error %q does not mention the offending version", err)
	}
}

// TestLoadV3Fixtures asserts every vendored 3.0 document loads into the 2.x
// model with populated core fields and well-formed relationship rows.
func TestLoadV3Fixtures(t *testing.T) {
	relTypeRe := regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	for _, name := range []string{
		"example3-bin.spdx3.json",
		"examplemaven-enriched.spdx3.json",
		"simplehtr-example.spdx3.json",
		"dataset-example01.spdx3.json",
		"rc4-compact-dialect.json",
	} {
		t.Run(name, func(t *testing.T) {
			s := loadFile(t, name)
			if s.SpdxVersion != "SPDX-3.0.1" {
				t.Errorf("SpdxVersion = %q, want SPDX-3.0.1", s.SpdxVersion)
			}
			if s.Spdxid == "" {
				t.Error("document id not mapped")
			}
			if len(s.CreationInfo.Creators) == 0 {
				t.Error("no creators mapped")
			}
			if s.CreationInfo.Created.IsZero() {
				t.Error("creation date not mapped")
			}
			// fixtures that are pure BOMs or datasets legitimately have
			// no Package elements; assert packages only where they exist
			if name == "example3-bin.spdx3.json" || name == "examplemaven-enriched.spdx3.json" {
				if len(s.Packages) == 0 {
					t.Error("no packages mapped")
				}
			}
			for _, rel := range s.Relationships {
				if rel.SpdxElementID == "" || rel.RelatedSpdxElement == "" {
					t.Errorf("relationship with empty endpoints: %+v", rel)
				}
				if !relTypeRe.MatchString(rel.RelationshipType) {
					t.Errorf("relationship type %q is not SCREAMING_SNAKE_CASE", rel.RelationshipType)
				}
			}
		})
	}
}

// TestLoadV3Example3 asserts concrete 2.x equivalences for the richest
// fixture: flattened lifecycle scopes, the generates -> GENERATED_FROM swap,
// hash and purl mapping.
func TestLoadV3Example3(t *testing.T) {
	s := loadFile(t, "example3-bin.spdx3.json")

	count := func(want string) int {
		n := 0
		for _, r := range s.Relationships {
			if r.RelationshipType == want {
				n++
			}
		}
		return n
	}
	if n := count("DEPENDS_ON"); n == 0 {
		t.Error("no DEPENDS_ON relationship rows (camelCase not converted?)")
	}
	if n := count("CONTAINS"); n == 0 {
		t.Error("no CONTAINS relationship rows")
	}
	// example3 models the BOM with 3.0 `generates` relationships; 2.x must
	// see GENERATED_FROM with the artifact as spdxElementId.
	if n := count("GENERATED_FROM"); n == 0 {
		t.Error("no GENERATED_FROM rows from 3.0 generates swap")
	}
	if n := count("HAS_CONCLUDED_LICENSE"); n == 0 {
		t.Error("licence relationships should still be emitted as rows")
	}
	if len(s.DocumentDescribes) == 0 {
		t.Error("documentDescribes not populated from describes relationships")
	}
	// rc4 needs injected placeholders for cross-document endpoints; they
	// must surface as ExternalElement rather than being lost.
	if n := s.SkippedElements["ExternalElement"]; n == 0 {
		t.Error("cross-document element references not recorded as ExternalElement")
	}
	// at least one package carries a licence through the 3.0
	// hasConcludedLicense/hasDeclaredLicense relationships
	var licensed bool
	for _, p := range s.Packages {
		if p.LicenseConcluded != "" || p.LicenseDeclared != "" {
			licensed = true
			break
		}
	}
	if !licensed {
		t.Error("no package carries a concluded/declared licence")
	}

	var file *Files
	for i := range s.Files {
		for _, c := range s.Files[i].Checksums {
			if c.Algorithm == "SHA256" && c.ChecksumValue != "" {
				file = &s.Files[i]
			}
		}
	}
	if file == nil {
		t.Fatal("no file with a SHA256 checksum")
	}
	var code string
	for _, p := range s.Packages {
		if c := p.PackageVerificationCode.PackageVerificationCodeValue; c != "" {
			code = c
		}
	}
	if code != "edbf31eb11b6b1698b7eec29bde0ea7040e0a9a4" {
		t.Errorf("package verification code = %q, want sha1 of main-bin", code)
	}
}

// TestLoadV3MavenPurl asserts packageUrl maps to a 2.x purl externalRef.
func TestLoadV3MavenPurl(t *testing.T) {
	s := loadFile(t, "examplemaven-enriched.spdx3.json")
	for _, p := range s.Packages {
		for _, ref := range p.ExternalRefs {
			if ref.ReferenceType == "purl" && strings.HasPrefix(ref.ReferenceLocator, "pkg:") {
				return
			}
		}
	}
	t.Error("no purl externalRef mapped")
}

// TestLoadV3SkippedElements asserts 3.0 elements without a 2.x mapping are
// counted for the meta sidecar row instead of silently dropped.
func TestLoadV3SkippedElements(t *testing.T) {
	s := loadFile(t, "simplehtr-example.spdx3.json")
	if got := s.SkippedElements["AIPackage"]; got != 2 {
		t.Errorf("SkippedElements[AIPackage] = %d, want 2", got)
	}
	if got := s.SkippedElements["DatasetPackage"]; got != 1 {
		t.Errorf("SkippedElements[DatasetPackage] = %d, want 1", got)
	}
	var core, software bool
	for _, p := range s.Profiles {
		switch p {
		case "core":
			core = true
		case "software":
			software = true
		}
	}
	if !core || !software {
		t.Errorf("Profiles = %v, want at least core and software", s.Profiles)
	}
	if len(s.HasExtractedLicensingInfos) == 0 {
		t.Error("custom licence text not mapped to hasExtractedLicensingInfos")
	}
	for _, e := range s.HasExtractedLicensingInfos {
		// 3.0 preserves the source licence id verbatim; the 2.x
		// LicenseRef- prefix convention is not enforced upstream
		if e.LicenseID == "" {
			t.Error("licence id not mapped")
		}
		if e.ExtractedText == "" {
			t.Errorf("licence %q has no extracted text", e.LicenseID)
		}
	}
}

// TestLoadV2Regression asserts plain 2.x JSON still parses with the 2.x path.
func TestLoadV2Regression(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "2x.json")
	body := `{
		"SPDXID": "SPDXRef-DOCUMENT",
		"spdxVersion": "SPDX-2.3",
		"name": "doc-2x",
		"creationInfo": {"created": "2024-01-01T00:00:00Z", "creators": ["Person: t"]},
		"packages": [{"SPDXID": "SPDXRef-Pkg", "name": "p", "versionInfo": "1.0"}],
		"relationships": [{"spdxElementId": "SPDXRef-DOCUMENT", "relationshipType": "DESCRIBES", "relatedSpdxElement": "SPDXRef-Pkg"}]
	}`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	s := new(SPDX)
	if err := s.Load(p); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.SpdxVersion != "SPDX-2.3" {
		t.Errorf("SpdxVersion = %q, want SPDX-2.3", s.SpdxVersion)
	}
	if len(s.Packages) != 1 || s.Packages[0].Name != "p" {
		t.Errorf("packages not parsed via 2.x path: %+v", s.Packages)
	}
	if len(s.Relationships) != 1 || s.Relationships[0].RelationshipType != "DESCRIBES" {
		t.Errorf("relationships not parsed via 2.x path: %+v", s.Relationships)
	}
	if s.Profiles != nil || len(s.SkippedElements) != 0 {
		t.Error("3.0 sidecar fields must stay empty on 2.x input")
	}
}
