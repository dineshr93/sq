// SPDX-FileCopyrightText: 2026 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"regexp"
	"strings"
	"testing"
)

// TestDigGeneratedFromResolvesNames probes the dig view against the
// spec-pinned 3.0.1 fixture (spdx-examples software/example3):
// "main-bin generates ./main" must show both sides by name, not blanks.
// Pre-fix the detail printer looked up only Packages, so file endpoints
// rendered as " | version: " with nothing around them.
func TestDigGeneratedFromResolvesNames(t *testing.T) {
	s := loadFile(t, "example3-bin.spdx3.json")
	s.GetRelationsforType()
	if len(s.RelTypes.GeneratedFrom) == 0 {
		t.Fatal("no GENERATED_FROM bucket")
	}
	for _, rel := range s.RelTypes.GeneratedFrom {
		from, _ := s.getPKGNameVersionDetailforRelsSPDXID(rel)
		plain := strings.TrimSpace(stripANSI(from))
		if plain == "" || strings.HasPrefix(plain, "|") || plain == rel.SpdxElementID {
			t.Errorf("GENERATED_FROM source unresolved: %q (%+v)", plain, rel)
		}
	}
}

var ansiRe = regexp.MustCompile("\x1b\\[[0-9;]*m")

func stripANSI(s string) string { return ansiRe.ReplaceAllString(s, "") }
