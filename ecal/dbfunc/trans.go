/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package dbfunc

import (
	"fmt"
	"strconv"

	"github.com/Fisch-Labs/FishDB/graph"
	"github.com/Fisch-Labs/ecal/parser"
)

/*
NewTransFunc creates a new transaction for FishDB.
*/
type NewTransFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *NewTransFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error

	if len(args) != 0 {
		err = fmt.Errorf("Function does not require any parameters")
	}

	return graph.NewConcurrentGraphTrans(f.GM), err
}

/*
DocString returns a descriptive string.
*/
func (f *NewTransFunc) DocString() (string, error) {
	return "Creates a new transaction for FishDB.", nil
}

/*
NewRollingTransFunc creates a new rolling transaction for FishDB.
A rolling transaction commits after n entries.
*/
type NewRollingTransFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *NewRollingTransFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error
	var trans graph.Trans

	if arglen := len(args); arglen != 1 {
		err = fmt.Errorf(
			"Function requires the rolling threshold (number of operations before rolling)")
	}

	if err == nil {
		var i int

		if i, err = strconv.Atoi(fmt.Sprint(args[0])); err != nil {
			err = fmt.Errorf("Rolling threshold must be a number not: %v", args[0])
		} else {
			trans = graph.NewRollingTrans(graph.NewConcurrentGraphTrans(f.GM),
				i, f.GM, graph.NewConcurrentGraphTrans)
		}
	}

	return trans, err
}

/*
DocString returns a descriptive string.
*/
func (f *NewRollingTransFunc) DocString() (string, error) {
	return "Creates a new rolling transaction for FishDB. A rolling transaction commits after n entries.", nil
}

/*
CommitTransFunc commits an existing transaction for FishDB.
*/
type CommitTransFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *CommitTransFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error

	if arglen := len(args); arglen != 1 {
		err = fmt.Errorf(
			"Function requires the transaction to commit as parameter")
	}

	if err == nil {
		trans, ok := args[0].(graph.Trans)

		// Check parameters

		if !ok {
			err = fmt.Errorf("Parameter must be a transaction")
		} else {
			err = trans.Commit()
		}
	}

	return nil, err
}

/*
DocString returns a descriptive string.
*/
func (f *CommitTransFunc) DocString() (string, error) {
	return "Commits an existing transaction for FishDB.", nil
}
