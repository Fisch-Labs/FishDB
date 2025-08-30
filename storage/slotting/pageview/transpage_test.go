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

func TestTransPage(t *testing.T) {
	r := file.NewRecord(123, make([]byte, 44))

	testCheckTransPageMagicPanic(t, r)

	// Make sure the record has a correct magic

	view.NewPageView(r, view.TypeTranslationPage)

	NewTransPage(r)
}

func testCheckTransPageMagicPanic(t *testing.T, r *file.Record) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Checking magic should fail.")
		}
	}()

	checkTransPageMagic(r)
}
