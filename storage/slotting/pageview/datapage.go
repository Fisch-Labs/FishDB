/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

/*
Package pageview contains object wrappers for different page types.

# DataPage

DataPage is a page which holds actual data.

# FreeLogicalSlotPage

FreeLogicalSlotPage is a page which holds information about free logical slots.
The page stores the slot location in a slotinfo data structure.

# FreePhysicalSlotPage

FreePhysicalSlotPage is a page which holds information about free physical slots.
The page stores the slot location and its size in a slotinfo data structure
(see util/slotsize.go).

# SlotInfoPage

SlotInfoPage is the super-struct for all page views which manage slotinfos.
Slotinfo are location (see util/location.go) pointers into the data store containing
record id and offset.

# TransPage

TransPage is a page which holds data to translate between physical and logical
slots.
*/
package pageview

import (
	"fmt"

	"github.com/Fisch-Labs/FishDB/storage/file"
	"github.com/Fisch-Labs/FishDB/storage/paging/view"
)

/*
OffsetFirst is a pointer to first element on the page
*/
const OffsetFirst = view.OffsetData

// OffsetData is declared in freephysicalslotpage

/*
DataPage data structure
*/
type DataPage struct {
	*SlotInfoPage
}

/*
NewDataPage creates a new page which holds actual data.
*/
func NewDataPage(record *file.Record) *DataPage {
	checkDataPageMagic(record)
	dp := &DataPage{NewSlotInfoPage(record)}
	return dp
}

/*
checkDataPageMagic checks if the magic number at the beginning of
the wrapped record is valid.
*/
func checkDataPageMagic(record *file.Record) bool {
	magic := record.ReadInt16(0)

	if magic == view.ViewPageHeader+view.TypeDataPage {
		return true
	}
	panic("Unexpected header found in DataPage")
}

/*
DataSpace returns the available data space on this page.
*/
func (dp *DataPage) DataSpace() uint16 {
	return uint16(len(dp.Record.Data()) - OffsetData)
}

/*
OffsetFirst returns the pointer to the first element on the page.
*/
func (dp *DataPage) OffsetFirst() uint16 {
	return dp.Record.ReadUInt16(OffsetFirst)
}

/*
SetOffsetFirst sets the pointer to the first element on the page.
*/
func (dp *DataPage) SetOffsetFirst(first uint16) {
	if first > 0 && first < OffsetData {
		panic(fmt.Sprint("Cannot set offset of first element on DataPage below ", OffsetData))
	}
	dp.Record.WriteUInt16(OffsetFirst, first)
}
