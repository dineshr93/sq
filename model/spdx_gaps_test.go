// SPDX-FileCopyrightText: 2026 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRelationsForTypeBuckets asserts GENERATED_FROM rows land in their own
// bucket, not CONTAINS, and that RelTypes.GeneratedFrom is populated.
func TestRelationsForTypeBuckets(t *testing.T) {
	s := loadFile(t, "example3-bin.spdx3.json")
	s.GetRelationsforType()
	for _, r := range s.RelTypes.GeneratedFrom {
		if r.RelationshipType != "GENERATED_FROM" {
			t.Errorf("GeneratedFrom bucket holds %q", r.RelationshipType)
		}
	}
	if len(s.RelTypes.GeneratedFrom) == 0 {
		t.Fatal("GeneratedFrom empty; GENERATED_FROM rows lost")
	}
	for _, r := range s.RelTypes.Contains {
		if r.RelationshipType == "GENERATED_FROM" {
			t.Errorf("GENERATED_FROM row wrongly in Contains bucket: %+v", r)
		}
	}
}

// TestRequestCountClamps asserts PrintFiles/Printpkgs stop at len when the
// user requests more rows than exist (sq pkgs 20 on a 5-package doc).
func TestRequestCountClamps(t *testing.T) {
	s := loadFile(t, "examplemaven-enriched.spdx3.json")
	if len(s.Packages) < 2 {
		t.Fatal("fixture too small")
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on oversized request: %v", r)
		}
	}()
	s.Printpkgs(len(s.Packages) + 10)
}

// TestLoadMissingFile asserts Load reports a missing file instead of
// silently producing an empty SPDX struct.
func TestLoadMissingFile(t *testing.T) {
	s := new(SPDX)
	if err := s.Load(filepath.Join(t.TempDir(), "nope.json")); err == nil {
		t.Error("Load(missing) returned nil error")
	}
}

// TestLoadEmptyFile asserts Load reports an empty file instead of nil error.
func TestLoadEmptyFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(p, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	s := new(SPDX)
	if err := s.Load(p); err == nil {
		t.Error("Load(empty) returned nil error")
	}
}
func TestDigDatasetPrintsFileRow(t *testing.T) {
	s := loadFile(t, "dataset-example01.spdx3.json")
	s.GetRelationsforType()
	var d = 0
	for _, rel := range s.RelTypes.Describes {
		d++
		if _, isPresent := s.getPKGforSPDXID(rel.SpdxElementID); isPresent {
			continue
		}
		if _, isPresentRel := s.getPKGforSPDXID(rel.RelatedSpdxElement); isPresentRel {
			continue
		}
		from := s.getPKGFileNameVersionforSPDXID(rel.SpdxElementID)
		to := s.getPKGFileNameVersionforSPDXID(rel.RelatedSpdxElement)
		if from == "" || to == "" {
			t.Errorf("file->file DESCRIBES row resolves empty endpoints: %+v", rel)
		}
	}
	if len(s.RelTypes.Describes) == 0 {
		t.Fatal("no DESCRIBES bucket")
	}
}
