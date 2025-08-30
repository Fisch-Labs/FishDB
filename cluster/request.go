/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package cluster

import "encoding/gob"

func init() {

	// Make sure we can use the relevant types in a gob operation

	gob.Register(&DataRequest{})
	gob.Register(make(map[string]string))
}

/*
RequestType is the type of a request
*/
type RequestType string

/*
List of all possible request types
*/
const (

	// Main DB

	RTGetMain RequestType = "GetMain"
	RTSetMain             = "SetMain"

	// Roots

	RTGetRoot = "GetRoot"
	RTSetRoot = "SetRoot"

	// Insert data

	RTInsert = "Insert"

	// Update data

	RTUpdate = "Update"

	// Free data

	RTFree = "Free"

	// Check for data

	RTExists = "Exists"

	// Retrieve data

	RTFetch = "Fetch"

	// Rebalance data

	RTRebalance = "Rebalance"
)

/*
DataRequestArg is a data request argument
*/
type DataRequestArg string

/*
List of all possible data request parameters.
*/
const (
	RPStoreName DataRequestArg = "StoreName" // Name of the store
	RPLoc                      = "Loc"       // Location of data
	RPVer                      = "Ver"       // Version of data
	RPRoot                     = "Root"      // Root id
	RPSrc                      = "Src"       // Request source member
)

/*
DataRequest data structure
*/
type DataRequest struct {
	RequestType RequestType                    // Type of request
	Args        map[DataRequestArg]interface{} // Request arguments
	Value       interface{}                    // Request value
	Transfer    bool                           // Flag for data transfer request
}
