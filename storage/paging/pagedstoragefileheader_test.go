/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package paging

import (
	"testing"

	"github.com/Fisch-Labs/FishDB/storage/file"
)

func TestPagedStorageFileHeader(t *testing.T) {

	record := file.NewRecord(5, make([]byte, 5, 5))
	testPagedStorageFileInitPanic1(t, record)

	record = file.NewRecord(5, make([]byte, 100, 100))
	testPagedStorageFileInitPanic2(t, record)

	NewPagedStorageFileHeader(record, true)
	psfh := NewPagedStorageFileHeader(record, true)

	if psfh.Roots() != 2 {
		t.Error("Unexpected number of roots:", psfh.Roots())
	}

	psfh.SetRoot(1, 0x42)
	if psfh.Root(1) != 0x42 {
		t.Error("Unexpected root value:", psfh.Root(1))
	}
	if psfh.Root(0) != 0x00 {
		t.Error("Unexpected root value:", psfh.Root(0))
	}

	psfh.SetFirstListElement(3, 5)
	if psfh.FirstListElement(3) != 5 {
		t.Error("Unexpected root value:", psfh.FirstListElement(3))
	}

	psfh.SetLastListElement(2, 5)
	if psfh.LastListElement(2) != 5 {
		t.Error("Unexpected root value:", psfh.LastListElement(3))
	}
}

func testPagedStorageFileInitPanic1(t *testing.T, r *file.Record) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Using a record which is too small did not cause a panic.")
		}
	}()

	NewPagedStorageFileHeader(r, true)
}

func testPagedStorageFileInitPanic2(t *testing.T, r *file.Record) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Using a record without header magic value did not cause a panic.")
		}
	}()

	NewPagedStorageFileHeader(r, false)
}
