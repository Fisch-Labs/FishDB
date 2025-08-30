/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package dbfunc

import (
	"fmt"

	"github.com/Fisch-Labs/FishDB/eql"
	"github.com/Fisch-Labs/FishDB/graph"
	"github.com/Fisch-Labs/ecal/parser"
)

/*
QueryFunc runs an EQL query.
*/
type QueryFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *QueryFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error
	var cols, rows []interface{}

	if arglen := len(args); arglen != 2 {
		err = fmt.Errorf("Function requires 2 parameters: partition and a query string")
	}

	if err == nil {
		var res eql.SearchResult

		part := fmt.Sprint(args[0])
		query := fmt.Sprint(args[1])

		res, err = eql.RunQuery("db.query", part, query, f.GM)

		if err != nil {
			return nil, err
		}

		// Convert result to rumble data structure

		labels := res.Header().Labels()
		cols = make([]interface{}, len(labels))
		for i, v := range labels {
			cols[i] = v
		}

		rrows := res.Rows()
		rows = make([]interface{}, len(rrows))
		for i, v := range rrows {
			rows[i] = v
		}
	}

	return map[interface{}]interface{}{
		"cols": cols,
		"rows": rows,
	}, err
}

/*
DocString returns a descriptive string.
*/
func (f *QueryFunc) DocString() (string, error) {
	return "Run an EQL query.", nil
}
