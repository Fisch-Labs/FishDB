/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package cluster

import (
	"fmt"
	"strconv"

	"github.com/Fisch-Labs/common/errorutil"
)

/*
toUInt64 safely converts an interface{} to an uint64.
*/
func toUInt64(v interface{}) uint64 {
	if vu, ok := v.(uint64); ok {
		return vu
	}

	cloc, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	errorutil.AssertOk(err)

	return uint64(cloc)
}
