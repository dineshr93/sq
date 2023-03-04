// SPDX-FileCopyrightText: 2023 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
)

type SPDX struct {
	Spdxid       string       `json:"SPDXID,omitempty"`
	Name         string       `json:"name,omitempty"`
	SpdxVersion  string       `json:"spdxVersion,omitempty"`
	CreationInfo CreationInfo `json:"creationInfo,omitempty"`

	DataLicense                string                       `json:"dataLicense,omitempty"`
	HasExtractedLicensingInfos []HasExtractedLicensingInfos `json:"hasExtractedLicensingInfos,omitempty"`
	DocumentNamespace          string                       `json:"documentNamespace,omitempty"`
	DocumentDescribes          []string                     `json:"documentDescribes,omitempty"`
	Files                      []Files                      `json:"files,omitempty"`
	Packages                   []Packages                   `json:"packages,omitempty"`

	Relationships []Relationships `json:"relationships,omitempty"`
	RelTypes      RelTypes        `json:"reltypes,omitempty"`
}
type CreationInfo struct {
	Created            time.Time `json:"created,omitempty"`
	Creators           []string  `json:"creators,omitempty"`
	LicenseListVersion string    `json:"licenseListVersion,omitempty"`
}
type Checksums struct {
	Algorithm     string `json:"algorithm,omitempty"`
	ChecksumValue string `json:"checksumValue,omitempty"`
}
type HasExtractedLicensingInfos struct {
	LicenseID     string `json:"licenseId,omitempty"`
	ExtractedText string `json:"extractedText,omitempty"`
	Name          string `json:"name,omitempty"`
}

type ExternalRefs struct {
	ReferenceCategory string `json:"referenceCategory,omitempty"`
	ReferenceLocator  string `json:"referenceLocator,omitempty"`
	ReferenceType     string `json:"referenceType,omitempty"`
}
type PackageVerificationCode struct {
	PackageVerificationCodeValue string `json:"packageVerificationCodeValue,omitempty"`
}
type Packages struct {
	Spdxid                string         `json:"SPDXID,omitempty"`
	Name                  string         `json:"name,omitempty"`
	FilesAnalyzed         bool           `json:"filesAnalyzed,omitempty"`
	CopyrightText         string         `json:"copyrightText,omitempty"`
	Description           string         `json:"description,omitempty"`
	DownloadLocation      string         `json:"downloadLocation,omitempty"`
	PrimaryPackagePurpose string         `json:"primaryPackagePurpose,omitempty"`
	Checksums             []Checksums    `json:"checksums,omitempty"`
	ExternalRefs          []ExternalRefs `json:"externalRefs,omitempty"`

	Homepage         string `json:"homepage,omitempty"`
	LicenseConcluded string `json:"licenseConcluded,omitempty"`
	LicenseDeclared  string `json:"licenseDeclared,omitempty"`

	Supplier                string                  `json:"supplier,omitempty"`
	Originator              string                  `json:"originator,omitempty"`
	VersionInfo             string                  `json:"versionInfo,omitempty"`
	SourceInfo              string                  `json:"sourceInfo,omitempty"`
	HasFiles                []string                `json:"hasFiles,omitempty"`
	PackageVerificationCode PackageVerificationCode `json:"packageVerificationCode,omitempty"`
}

type Files struct {
	Spdxid             string      `json:"SPDXID,omitempty"`
	Checksums          []Checksums `json:"checksums,omitempty"`
	CopyrightText      string      `json:"copyrightText,omitempty"`
	FileName           string      `json:"fileName,omitempty"`
	LicenseConcluded   string      `json:"licenseConcluded,omitempty"`
	LicenseInfoInFiles []string    `json:"licenseInfoInFiles,omitempty"`
}
type Relationships struct {
	SpdxElementID      string `json:"spdxElementId,omitempty"`
	RelationshipType   string `json:"relationshipType,omitempty"`
	RelatedSpdxElement string `json:"relatedSpdxElement,omitempty"`
}
type RelTypes struct {
	Describes     []Relationships
	Dependson     []Relationships
	Contains      []Relationships
	GeneratedFrom []Relationships
}

func (s *SPDX) GetRelationsforType() {
	rels := s.Relationships

	var relDescribes []Relationships
	var relDependson []Relationships
	var relContains []Relationships
	var generatedFrom []Relationships
	for _, rel := range rels {
		switch rt := rel.RelationshipType; rt {
		case "DESCRIBES":
			relDescribes = append(relDescribes, rel)
		case "DEPENDS_ON":
			relDependson = append(relDependson, rel)
		case "CONTAINS":
			relContains = append(relContains, rel)
		case "GENERATED_FROM":
			relContains = append(generatedFrom, rel)
		}
	}
	relTypes := RelTypes{
		Describes: relDescribes,
		Dependson: relDependson,
		Contains:  relContains,
	}
	s.RelTypes = relTypes
}
func (t *SPDX) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (s *SPDX) PrintMeta() {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Key"},
			{Align: simpletable.AlignCenter, Text: "Value"},
		},
	}

	var cells [][]*simpletable.Cell

	idx := 0

	idx++

	cells = append(cells, []*simpletable.Cell{

		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Spdx ID")},
		{Text: s.Spdxid},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Spdx version")},
		{Text: s.SpdxVersion},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Spdx creation date")},
		{Text: yellow(s.CreationInfo.Created.Format(time.RFC822))},
	})
	idx++
	if len(s.CreationInfo.Creators) > 0 {
		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: blue("created by")},
			{Text: strings.Join(s.CreationInfo.Creators, ", ")},
		})
	}
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Project Name")},
		{Text: red(s.Name)},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("File License(not projects)")},
		{Text: s.DataLicense},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Document Namespace")},
		{Text: s.DocumentNamespace},
	})
	idx++
	var pkgVersionNames []string
	for _, spdxID := range s.DocumentDescribes {
		pkg, isPresent := s.getPKGforSPDXID(spdxID)
		if isPresent {
			pkgVersionNames = append(pkgVersionNames, fmt.Sprintf("%v-%v", pkg.Name, pkg.VersionInfo))
		}
	}
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Document Describes")},

		{Text: strings.Join(pkgVersionNames, ", ")},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Packages")},
		{Text: red(fmt.Sprintf("%d", len(s.Packages)))},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Files")},
		{Text: red(fmt.Sprintf("%d", len(s.Files)))},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Relationships")},
		{Text: red(fmt.Sprintf("%d", len(s.Relationships)))},
	})

	table.Body = &simpletable.Body{Cells: cells}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}

func (s *SPDX) PrintFiles(nf int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "List of Files"},
			// {Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			// {Align: simpletable.AlignCenter, Text: "LicenseInfoInFiles"},
			// {Align: simpletable.AlignCenter, Text: "SPDXId"},
			// {Align: simpletable.AlignCenter, Text: "CopyrightText"},
			// {Align: simpletable.AlignCenter, Text: "checksum"},
			// {Align: simpletable.AlignCenter, Text: "Algorithm"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	files := s.Files
	var file Files
	var n int
	// var licenseinfo string
	lenFiles := len(files)
	// fmt.Println(lenFiles)
	var cells [][]*simpletable.Cell

	for id := 0; id < nf; id++ {

		file = files[id]

		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: red(file.FileName)},
			// {Text: file.LicenseConcluded},
			// {Text: licenseinfo},
			// {Text: file.Spdxid},
			// {Text: file.CopyrightText},
			// {Text: file.Checksums[0].ChecksumValue},
			// {Text: file.Checksums[0].Algorithm},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenFiles > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 2, Text: blue(fmt.Sprintf("There are %d Files", lenFiles))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (s *SPDX) PrintFilesCheksum(nf int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "FileName"},
			// {Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			// {Align: simpletable.AlignCenter, Text: "LicenseInfoInFiles"},
			// {Align: simpletable.AlignCenter, Text: "SPDXId"},
			// {Align: simpletable.AlignCenter, Text: "CopyrightText"},
			{Align: simpletable.AlignCenter, Text: "checksum"},
			{Align: simpletable.AlignCenter, Text: "Algorithm"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	files := s.Files
	var file Files
	var n int
	// var licenseinfo string
	lenFiles := len(files)
	// fmt.Println(lenFiles)
	var cells [][]*simpletable.Cell

	for id := 0; id < nf; id++ {

		file = files[id]

		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: file.FileName},
			// {Text: file.LicenseConcluded},
			// {Text: licenseinfo},
			// {Text: file.Spdxid},
			// {Text: file.CopyrightText},
			{Text: file.Checksums[0].ChecksumValue},
			{Text: file.Checksums[0].Algorithm},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenFiles > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 4, Text: blue(fmt.Sprintf("There are %d Files", lenFiles))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (s *SPDX) PrintFilesIP(nf int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "FileName"},
			{Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			{Align: simpletable.AlignCenter, Text: "LicenseInfoInFiles"},
			// {Align: simpletable.AlignCenter, Text: "SPDXId"},
			{Align: simpletable.AlignCenter, Text: "CopyrightText"},
			// {Align: simpletable.AlignCenter, Text: "checksum"},
			// {Align: simpletable.AlignCenter, Text: "Algorithm"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	files := s.Files
	var file Files
	var n int
	var licenseinfo string
	lenFiles := len(files)
	// fmt.Println(lenFiles)
	var cells [][]*simpletable.Cell

	for id := 0; id < nf; id++ {

		file = files[id]

		n = id + 1
		licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: file.FileName},
			{Text: file.LicenseConcluded},
			{Text: licenseinfo},
			// {Text: file.Spdxid},
			{Text: file.CopyrightText},
			// {Text: file.Checksums[0].ChecksumValue},
			// {Text: file.Checksums[0].Algorithm},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenFiles > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: blue(fmt.Sprintf("There are %d Files", lenFiles))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}
func (s *SPDX) PrintFilesExt(nf int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "FileName"},
			{Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			{Align: simpletable.AlignCenter, Text: "LicenseInfoInFiles"},
			{Align: simpletable.AlignCenter, Text: "SPDXId"},
			{Align: simpletable.AlignCenter, Text: "CopyrightText"},
			{Align: simpletable.AlignCenter, Text: "checksum"},
			{Align: simpletable.AlignCenter, Text: "Algorithm"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	files := s.Files
	var file Files
	var n int
	var licenseinfo string
	lenFiles := len(files)
	// fmt.Println(lenFiles)
	var cells [][]*simpletable.Cell

	for id := 0; id < nf; id++ {

		file = files[id]

		n = id + 1
		licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: file.FileName},
			{Text: file.LicenseConcluded},
			{Text: licenseinfo},
			{Text: file.Spdxid},
			{Text: file.CopyrightText},
			{Text: file.Checksums[0].ChecksumValue},
			{Text: file.Checksums[0].Algorithm},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenFiles > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 8, Text: blue(fmt.Sprintf("There are %d Files", lenFiles))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func getPkgNameVersion(pkg Packages) (string, string) {
	var isSha bool
	var pkgName string = pkg.Name
	var pkgVersion string = pkg.VersionInfo
	for _, d := range shas {
		if strings.HasPrefix(strings.ToLower(pkgName), d) {
			isSha = true
			break
		}
	}
	if isSha {
		pkgName = pkg.Description
		pkgVersion = pkg.PrimaryPackagePurpose
	}
	return pkgName, pkgVersion
}

func (s *SPDX) Printpkgs(np int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Supplier"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "VersionInfo"},
			{Align: simpletable.AlignCenter, Text: "Homepage"},
			// {Align: simpletable.AlignCenter, Text: "LicenseDeclared"},
			// {Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			// {Align: simpletable.AlignCenter, Text: "FilesAnalyzed"},
			// {Align: simpletable.AlignCenter, Text: "DownloadLocation"},
			// {Align: simpletable.AlignCenter, Text: "CopyrightText"},
			// {Align: simpletable.AlignCenter, Text: "Spdxid"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	pkgs := s.Packages
	var pkg Packages
	var n int
	// var licenseinfo string
	lenPkgs := len(pkgs)
	// fmt.Println(lenPkgs)
	var cells [][]*simpletable.Cell

	for id := 0; id < np; id++ {

		pkg = pkgs[id]
		// ===============================
		pkgName, pkgVersion := getPkgNameVersion(pkg)
		// ===================================

		homepage := pkg.Homepage
		if homepage == "" {
			homepage = pkg.DownloadLocation
		}
		supplier := pkg.Supplier
		if supplier == "" {
			supplier = pkg.Originator
		}
		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: supplier},
			{Text: pkgName},
			{Text: pkgVersion},
			{Text: homepage},
			// {Text: pkg.LicenseDeclared},
			// {Text: pkg.LicenseConcluded},
			// {Text: fmt.Sprintf("%v", pkg.FilesAnalyzed)},
			// {Text: pkg.DownloadLocation},
			// {Text: pkg.CopyrightText},
			// {Text: pkg.Spdxid},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenPkgs > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 5, Text: blue(fmt.Sprintf("There are %d pkgs", lenPkgs))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

}

func (s *SPDX) PrintpkgsIP(np int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			// {Align: simpletable.AlignCenter, Text: "Supplier"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "VersionInfo"},
			// {Align: simpletable.AlignCenter, Text: "Homepage"},
			{Align: simpletable.AlignCenter, Text: "LicenseDeclared"},
			{Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			{Align: simpletable.AlignCenter, Text: "FilesAnalyzed"},
			// {Align: simpletable.AlignCenter, Text: "DownloadLocation"},
			{Align: simpletable.AlignCenter, Text: "CopyrightText"},
			// {Align: simpletable.AlignCenter, Text: "Spdxid"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	pkgs := s.Packages
	var pkg Packages
	var n int
	// var licenseinfo string
	lenPkgs := len(pkgs)
	// fmt.Println(lenPkgs)
	var cells [][]*simpletable.Cell

	for id := 0; id < np; id++ {

		pkg = pkgs[id]

		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			// {Text: pkg.Supplier},
			{Text: pkg.Name},
			{Text: pkg.VersionInfo},
			// {Text: pkg.Homepage},
			{Text: pkg.LicenseDeclared},
			{Text: pkg.LicenseConcluded},
			{Text: fmt.Sprintf("%v", pkg.FilesAnalyzed)},
			// {Text: pkg.DownloadLocation},
			{Text: pkg.CopyrightText},

			// {Text: pkg.Spdxid},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenPkgs > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 7, Text: blue(fmt.Sprintf("There are %d pkgs", lenPkgs))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

}
func (s *SPDX) PrintpkgsExt(np int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Supplier"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "VersionInfo"},
			{Align: simpletable.AlignCenter, Text: "Homepage"},
			{Align: simpletable.AlignCenter, Text: "LicenseDeclared"},
			{Align: simpletable.AlignCenter, Text: "LicenseConcluded"},
			{Align: simpletable.AlignCenter, Text: "FilesAnalyzed"},
			{Align: simpletable.AlignCenter, Text: "DownloadLocation"},
			{Align: simpletable.AlignCenter, Text: "CopyrightText"},
			// {Align: simpletable.AlignCenter, Text: "Spdxid"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	pkgs := s.Packages
	var pkg Packages
	var n int
	// var licenseinfo string
	lenPkgs := len(pkgs)
	// fmt.Println(lenPkgs)
	var cells [][]*simpletable.Cell

	for id := 0; id < np; id++ {

		pkg = pkgs[id]

		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", n)},
			{Text: pkg.Supplier},
			{Text: pkg.Name},
			{Text: pkg.VersionInfo},
			{Text: pkg.Homepage},
			{Text: pkg.LicenseDeclared},
			{Text: pkg.LicenseConcluded},
			{Text: fmt.Sprintf("%v", pkg.FilesAnalyzed)},
			{Text: pkg.DownloadLocation},
			{Text: pkg.CopyrightText},
			// {Text: pkg.Spdxid},
		})

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenPkgs > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 10, Text: blue(fmt.Sprintf("There are %d pkgs", lenPkgs))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

}

func (s *SPDX) PrintRelsinSPDX(np int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "SpdxElementID"},
			{Align: simpletable.AlignCenter, Text: "RelationshipType"},
			{Align: simpletable.AlignCenter, Text: "RelatedSpdxElement"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	rels := s.Relationships
	var rel Relationships
	var n int
	// var licenseinfo string
	lenrels := len(rels)
	// fmt.Println(lenrels)
	var cells [][]*simpletable.Cell

	for id := 0; id < np; id++ {

		rel = rels[id]

		n = id + 1
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")
		switch rt := rel.RelationshipType; rt {
		case "CONTAINS":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: rel.SpdxElementID},
				{Text: red(rt)},
				{Text: rel.RelatedSpdxElement},
			})
		case "DEPENDS_ON":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: rel.SpdxElementID},
				{Text: blue(rt)},
				{Text: rel.RelatedSpdxElement},
			})
		case "DESCRIBES":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: rel.SpdxElementID},
				{Text: green(rt)},
				{Text: rel.RelatedSpdxElement},
			})
		default:
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: rel.SpdxElementID},
				{Text: gray(rt)},
				{Text: rel.RelatedSpdxElement},
			})

		}

	}
	table.Body = &simpletable.Body{Cells: cells}

	if lenrels > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 4, Text: blue(fmt.Sprintf("There are %d relationships", lenrels))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

}

func (s *SPDX) PrintRelsClarified(np int) {

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "SpdxElement"},
			{Align: simpletable.AlignCenter, Text: "RelationshipType"},
			{Align: simpletable.AlignCenter, Text: "RelatedSpdxElement"},
		},
	}
	// {Align: simpletable.AlignCenter, Text: "Checksums"},
	// {Align: simpletable.AlignCenter, Text: "Algorithm - Checksums"},
	rels := s.Relationships
	var rel Relationships
	var n int
	// var licenseinfo string
	lenrels := len(rels)
	// fmt.Println(lenrels)
	var cells [][]*simpletable.Cell
	var relContainsCount int

	for id := 0; id < np; id++ {

		rel = rels[id]

		n = id + 1
		SpdxElementID := s.getPKGFileNameVersionforSPDXID(rel.SpdxElementID)
		RelatedSpdxElement := s.getPKGFileNameVersionforSPDXID(rel.RelatedSpdxElement)
		// licenseinfo = strings.Join(file.LicenseInfoInFiles, ", ")
		switch rt := rel.RelationshipType; rt {
		case "CONTAINS":
			relContainsCount++
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: SpdxElementID},
				{Text: red(rt)},
				{Text: red(RelatedSpdxElement)},
			})
		case "DEPENDS_ON":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: SpdxElementID},
				{Text: blue(rt)},
				{Text: blue(RelatedSpdxElement)},
			})
		case "DESCRIBES":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: SpdxElementID},
				{Text: green(rt)},
				{Text: green(RelatedSpdxElement)},
			})
		case "GENERATED_FROM":
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: SpdxElementID},
				{Text: yellow(rt)},
				{Text: yellow(RelatedSpdxElement)},
			})
		default:
			cells = append(cells, *&[]*simpletable.Cell{
				{Text: fmt.Sprintf("%d", n)},
				{Text: SpdxElementID},
				{Text: gray(rt)},
				{Text: RelatedSpdxElement},
			})

		}

	}

	table.Body = &simpletable.Body{Cells: cells}

	if lenrels > 0 {
		table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Span: 4, Text: blue(fmt.Sprintf("There are %d relationships", lenrels))},
		}}
	}
	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

	if relContainsCount == 0 {
		s.PrintFiles(len(s.Files))
	}

}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

var shas []string = []string{"sha1", "sha256", "sha512"}

func (s *SPDX) getPKGFileNameVersionforSPDXID(spdxid string) string {
	pkg, ispkgPresent := s.getPKGforSPDXID(spdxid)
	if ispkgPresent {
		pkgName, pkgVersion := getPkgNameVersion(pkg)
		return fmt.Sprintf("%v %v", pkgName, pkgVersion)
	} else {
		file, isFilePresent := s.getFileforSPDXID(spdxid)
		if isFilePresent {
			return file.FileName
		}
	}
	return spdxid
}
func (s *SPDX) getPKGFileforSPDXID(spdxid string) (Packages, bool, Files, bool) {
	pkg, ispkgPresent := s.getPKGforSPDXID(spdxid)
	if ispkgPresent {
		return pkg, true, Files{}, false
	} else {
		file, isFilePresent := s.getFileforSPDXID(spdxid)

		if isFilePresent {
			return Packages{}, false, file, true
		}
	}
	return Packages{}, false, Files{}, false
}

func (s *SPDX) getPKGforSPDXID(spdxid string) (Packages, bool) {
	pkgs := s.Packages
	for _, pkg := range pkgs {
		if spdxid == pkg.Spdxid {
			return pkg, true
		}
	}
	return Packages{}, false
}
func (s *SPDX) getFileforSPDXID(spdxid string) (Files, bool) {
	files := s.Files
	for _, file_ := range files {
		if spdxid == file_.Spdxid {
			return file_, true
		}
	}
	return Files{}, false
}
func (s *SPDX) getRelationshipForSPDXID(spdxid string, rels []Relationships) []Relationships {
	var relTmp []Relationships
	for _, rel := range rels {
		if spdxid == rel.SpdxElementID {
			relTmp = append(relTmp, rel)
		}
	}
	return relTmp
}
func (s *SPDX) getPKGNameVersionDetailforRelsSPDXID(rel Relationships) (string, string) {
	var SPDXidDetail, relatedSPDXidDetail string
	spdxpkg, _ := s.getPKGforSPDXID(rel.SpdxElementID)
	SPDXidDetail = fmt.Sprintf("%v %v version: %v", blue(spdxpkg.Name), yellow("|"), blue(spdxpkg.VersionInfo))
	relatedspdxpkg, _ := s.getPKGforSPDXID(rel.RelatedSpdxElement)
	if len(relatedspdxpkg.HasFiles) > 0 {
		relatedSPDXidDetail = fmt.Sprintf("%v %v version: %v", yellow(relatedspdxpkg.Name), yellow("|"), yellow(relatedspdxpkg.VersionInfo))
	} else {
		relatedSPDXidDetail = fmt.Sprintf("%v %v version: %v", blue(relatedspdxpkg.Name), yellow("|"), blue(relatedspdxpkg.VersionInfo))
	}
	return SPDXidDetail, relatedSPDXidDetail
}

func (s *SPDX) printspdxpkg(d int, spdxID string, rel Relationships) {
	var filespdxids []string

	pkg, _ := s.getPKGforSPDXID(spdxID)
	// var SPDXidDetail string
	// SPDXidDetail = fmt.Sprintf("%v | version: %v", yellow(pkg.Name), yellow(pkg.VersionInfo))

	if len(pkg.HasFiles) > 0 {
		filespdxids = pkg.HasFiles
		// fmt.Println(red("    |-----------Digging one more level deep-----------"))
		// fmt.Println(green(fmt.Sprintf("    |-->%v %v %v", SPDXidDetail, yellow("---->"), red("CONTAINS"))))
		for i, filespdx := range filespdxids {
			s.printspdxfile(i, filespdx)
		}
		// fmt.Println()
	} else {
		// To Dig for files for pkg which doesnt have parameter HasFiles
		// All relationships with particular spdxid with Contains relationship
		// Dig for SpdxElementID
		rels := s.getRelationshipForSPDXID(rel.SpdxElementID, s.RelTypes.Contains)
		for f, rel := range rels {
			fmt.Println()
			// fmt.Println(red("    |-----------Digging one more level deep-----------"))
			// fmt.Println(green(fmt.Sprintf("    |-->%v %v %v", SPDXidDetail, yellow("---->"), red("CONTAINS"))))
			s.printspdxfile(f, rel.RelatedSpdxElement)
			// fmt.Println()
		}
		// Dig for RelatedSpdxElement
		rels = s.getRelationshipForSPDXID(rel.RelatedSpdxElement, s.RelTypes.Contains)
		for f, rel := range rels {
			// fmt.Println()
			// fmt.Println(red("    |-----------Digging one more level deep-----------"))
			// fmt.Println(green(fmt.Sprintf("    |-->%v %v %v", SPDXidDetail, yellow("---->"), red("CONTAINS"))))
			s.printspdxfile(f, rel.RelatedSpdxElement)
			// fmt.Println()
		}
	}
}

func (s *SPDX) printspdxfile(i int, spdxID string) {
	file, _ := s.getFileforSPDXID(spdxID)
	i++
	fmt.Println(green(fmt.Sprintf("    |-->File %v %v %v", blue(fmt.Sprintf("%d", i)), yellow("---->"), red(file.FileName))))
}
func (s *SPDX) printdependson(d int, rel Relationships) {
	_, relatedSPDXidDetail := s.getPKGNameVersionDetailforRelsSPDXID(rel)
	fmt.Println(fmt.Sprintf(" |-->Pkg %v %v", blue(fmt.Sprintf("%d", d)), relatedSPDXidDetail))
	// fmt.Println(fmt.Sprintf("%v ====> %v Pkg %v =====> %v", SPDXidDetail, green("DEPENDS_ON"), blue(fmt.Sprintf("%d", d)), relatedSPDXidDetail))
}
func (s *SPDX) PrintDigRels() {
	// Load struct RelTypes based on types
	s.GetRelationsforType()
	var tmp string

	for d, rel := range s.RelTypes.Describes {
		d++
		pkg, isPresent := s.getPKGforSPDXID(rel.SpdxElementID)
		pkgRel, isPresentRel := s.getPKGforSPDXID(rel.RelatedSpdxElement)
		if isPresent {
			fmt.Println(green("===================DESCRIBES/CONTAINS======================"))
			fmt.Println(fmt.Sprintf("Root Element %v %v %v %v %v", yellow(pkg.Name), yellow(pkg.VersionInfo), green("DESCRIBES"), yellow(pkgRel.Name), yellow(pkgRel.VersionInfo)))
		} else if isPresentRel {
			fmt.Println(green("===================DESCRIBES/CONTAINS======================"))
			fmt.Println(fmt.Sprintf("Root Element %v %v %v %v", yellow(rel.SpdxElementID), green("DESCRIBES"), yellow(pkgRel.Name), yellow(pkgRel.VersionInfo)))
		}
		s.printspdxpkg(d, rel.RelatedSpdxElement, rel)
	}

	if len(s.RelTypes.GeneratedFrom) > 0 {
		fmt.Println(blue("===================GENERATED_FROM======================"))
		tmp = ""
		for d, rel := range s.RelTypes.GeneratedFrom {
			d++
			SPDXidDetail, _ := s.getPKGNameVersionDetailforRelsSPDXID(rel)
			if tmp != SPDXidDetail {
				fmt.Println(fmt.Sprintf("%v %v", SPDXidDetail, blue("GENERATED_FROM")))
			}
			// s.printdependson(d, rel)
			s.printspdxpkg(d, rel.RelatedSpdxElement, rel)
			tmp = SPDXidDetail
		}
	}
	if len(s.RelTypes.Dependson) > 0 {
		fmt.Println(blue("===================DEPENDS_ON/CONTAINS======================"))
		tmp = ""
		for d, rel := range s.RelTypes.Dependson {
			d++
			SPDXidDetail, _ := s.getPKGNameVersionDetailforRelsSPDXID(rel)
			if tmp != SPDXidDetail {
				fmt.Println(fmt.Sprintf("%v %v", SPDXidDetail, blue("DEPENDS_ON")))
			}
			s.printdependson(d, rel)
			s.printspdxpkg(d, rel.RelatedSpdxElement, rel)
			tmp = SPDXidDetail
		}
	}

}
