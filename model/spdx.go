package model

import "time"

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
