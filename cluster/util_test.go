/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package cluster

import "testing"

func TestToUInt64(t *testing.T) {

	if res := toUInt64("1"); res != 1 {
		t.Error("Unexpected result: ", res)
		return
	}

	if res := toUInt64(uint64(1)); res != 1 {
		t.Error("Unexpected result: ", res)
		return
	}
}
