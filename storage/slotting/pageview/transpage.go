/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package pageview

import (
	"github.com/Fisch-Labs/FishDB/storage/file"
	"github.com/Fisch-Labs/FishDB/storage/paging/view"
)

/*
OffsetTransData is the data offset for translation pages
*/
const OffsetTransData = view.OffsetData

/*
TransPage data structure
*/
type TransPage struct {
	*SlotInfoPage
}

/*
NewTransPage creates a new page which holds data to translate between physical
and logical slots.
*/
func NewTransPage(record *file.Record) *DataPage {
	checkTransPageMagic(record)
	return &DataPage{NewSlotInfoPage(record)}
}

/*
checkTransPageMagic checks if the magic number at the beginning of
the wrapped record is valid.
*/
func checkTransPageMagic(record *file.Record) bool {
	magic := record.ReadInt16(0)

	if magic == view.ViewPageHeader+view.TypeTranslationPage {
		return true
	}
	panic("Unexpected header found in TransPage")
}
