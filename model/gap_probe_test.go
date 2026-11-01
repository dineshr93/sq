// SPDX-FileCopyrightText: 2026 Dinesh Ravi
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"testing"
)

// TestChecksumNoChecksumFile asserts PrintFilesCheksum does not panic on a
// file without checksums (real SBOMs contain such files) and prints
// NOASSERTION instead.
func TestChecksumNoChecksumFile(t *testing.T) {
	s := new(SPDX)
	s.Files = []Files{{FileName: "nochecksum.c", Spdxid: "SPDXRef-F1"}}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic on file without checksum: %v", r)
		}
	}()
	s.PrintFilesCheksum(1)
	if got := firstChecksum(s.Files[0]); got != noassert {
		t.Errorf("firstChecksum(no checksum) = %q, want %q", got, noassert)
	}
	if got := firstChecksumAlg(s.Files[0]); got != noassert {
		t.Errorf("firstChecksumAlg(no checksum) = %q, want %q", got, noassert)
	}
	withSum := Files{FileName: "a.c", Checksums: []Checksums{{Algorithm: "SHA1", ChecksumValue: "abc"}}}
	if got := firstChecksum(withSum); got != "abc" {
		t.Errorf("firstChecksum = %q, want abc", got)
	}
}
