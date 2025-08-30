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
	"github.com/Fisch-Labs/FishDB/storage/util"
)

/*
SlotInfoPage data structure
*/
type SlotInfoPage struct {
	*view.PageView
}

/*
NewSlotInfoPage creates a new SlotInfoPage object which can manage slotinfos.
*/
func NewSlotInfoPage(record *file.Record) *SlotInfoPage {
	pv := view.GetPageView(record)
	return &SlotInfoPage{pv}
}

/*
SlotInfoRecord gets record id of a stored slotinfo.
*/
func (lm *SlotInfoPage) SlotInfoRecord(offset uint16) uint64 {
	return util.LocationRecord(lm.Record.ReadUInt64(int(offset)))
}

/*
SlotInfoOffset gets the record offset of a stored slotinfo.
*/
func (lm *SlotInfoPage) SlotInfoOffset(offset uint16) uint16 {
	return util.LocationOffset(lm.Record.ReadUInt64(int(offset)))
}

/*
SetSlotInfo stores a slotinfo on the pageview's record.
*/
func (lm *SlotInfoPage) SetSlotInfo(slotinfoOffset uint16, recordID uint64, offset uint16) {
	lm.Record.WriteUInt64(int(slotinfoOffset), util.PackLocation(recordID, offset))
}
