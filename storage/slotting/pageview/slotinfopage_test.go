/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package pageview

import (
	"testing"

	"github.com/Fisch-Labs/FishDB/storage/file"
	"github.com/Fisch-Labs/FishDB/storage/paging/view"
)

func TestSlotInfoPage(t *testing.T) {
	r := file.NewRecord(123, make([]byte, 20))

	// Make sure the record has a correct magic

	view.NewPageView(r, view.TypeDataPage)

	si := NewSlotInfoPage(r)

	si.SetSlotInfo(2, 99, 45)

	if si.SlotInfoOffset(2) != 45 {
		t.Error("Unexpected offset read back")
	}

	if si.SlotInfoRecord(2) != 99 {
		t.Error("Unexpected record read back")
	}
}
