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

func TestDataPage(t *testing.T) {
	r := file.NewRecord(123, make([]byte, 44))

	testCheckDataPageMagicPanic(t, r)

	// Make sure the record has a correct magic

	view.NewPageView(r, view.TypeDataPage)

	dp := NewDataPage(r)

	if ds := dp.DataSpace(); ds != 24 {
		t.Error("Unexpected data space", ds)
	}

	testCheckDataPageOffsetFirstPanic(t, dp)

	if of := dp.OffsetFirst(); of != 0 {
		t.Error("Unexpected first offset", of)
		return
	}

	dp.SetOffsetFirst(20)

	if of := dp.OffsetFirst(); of != 20 {
		t.Error("Unexpected first offset", of)
		return
	}
}

func testCheckDataPageMagicPanic(t *testing.T, r *file.Record) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Checking magic should fail.")
		}
	}()

	checkDataPageMagic(r)
}

func testCheckDataPageOffsetFirstPanic(t *testing.T, dp *DataPage) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Setting offset which is too small should fail.")
		}
	}()

	dp.SetOffsetFirst(OffsetData - 1)
}
