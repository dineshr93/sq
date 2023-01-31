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
	Spdxid                     string                       `json:"SPDXID,omitempty"`
	SpdxVersion                string                       `json:"spdxVersion,omitempty"`
	CreationInfo               CreationInfo                 `json:"creationInfo,omitempty"`
	Name                       string                       `json:"name,omitempty"`
	DataLicense                string                       `json:"dataLicense,omitempty"`
	HasExtractedLicensingInfos []HasExtractedLicensingInfos `json:"hasExtractedLicensingInfos,omitempty"`
	DocumentNamespace          string                       `json:"documentNamespace,omitempty"`
	DocumentDescribes          []string                     `json:"documentDescribes,omitempty"`
	Packages                   []Packages                   `json:"packages,omitempty"`
	Files                      []Files                      `json:"files,omitempty"`
	Relationships              []Relationships              `json:"relationships,omitempty"`
}
type CreationInfo struct {
	Created            time.Time `json:"created,omitempty"`
	Creators           []string  `json:"creators,omitempty"`
	LicenseListVersion string    `json:"licenseListVersion,omitempty"`
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
type Packages struct {
	Spdxid           string         `json:"SPDXID,omitempty"`
	CopyrightText    string         `json:"copyrightText,omitempty"`
	DownloadLocation string         `json:"downloadLocation,omitempty"`
	ExternalRefs     []ExternalRefs `json:"externalRefs,omitempty"`
	FilesAnalyzed    bool           `json:"filesAnalyzed,omitempty"`
	Homepage         string         `json:"homepage,omitempty"`
	LicenseConcluded string         `json:"licenseConcluded,omitempty"`
	LicenseDeclared  string         `json:"licenseDeclared,omitempty"`
	Name             string         `json:"name,omitempty"`
	Supplier         string         `json:"supplier,omitempty"`
	VersionInfo      string         `json:"versionInfo,omitempty"`
	HasFiles         []string       `json:"hasFiles,omitempty"`
}
type Checksums struct {
	Algorithm     string `json:"algorithm,omitempty"`
	ChecksumValue string `json:"checksumValue,omitempty"`
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
	RelatedSpdxElement string `json:"relatedSpdxElement,omitempty"`
	RelationshipType   string `json:"relationshipType,omitempty"`
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
		{Text: s.CreationInfo.Created.Format(time.RFC822)},
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
		{Text: blue("Project version")},
		{Text: s.Name},
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
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Document Describes")},
		{Text: strings.Join(s.DocumentDescribes, ", ")},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Packages")},
		{Text: fmt.Sprintf("%d", len(s.Packages))},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Files")},
		{Text: fmt.Sprintf("%d", len(s.Files))},
	})
	idx++
	cells = append(cells, []*simpletable.Cell{
		{Text: fmt.Sprintf("%d", idx)},
		{Text: blue("Number of Relationships")},
		{Text: fmt.Sprintf("%d", len(s.Relationships))},
	})

	table.Body = &simpletable.Body{Cells: cells}

	table.SetStyle(simpletable.StyleUnicode)

	table.Println()
}
