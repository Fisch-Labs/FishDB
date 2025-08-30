/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

/*
Package dbfunc contains FishDB specific functions for the event condition action language (ECAL).
*/
package dbfunc

import (
	"fmt"

	"github.com/Fisch-Labs/FishDB/graph"
	"github.com/Fisch-Labs/FishDB/graph/data"
	"github.com/Fisch-Labs/ecal/parser"
)

/*
StoreEdgeFunc inserts or updates an edge in FishDB.
*/
type StoreEdgeFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *StoreEdgeFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error

	if arglen := len(args); arglen != 2 && arglen != 3 {
		err = fmt.Errorf("Function requires 2 or 3 parameters: partition, edge" +
			" map and optionally a transaction")
	}

	if err == nil {
		var trans graph.Trans

		part := fmt.Sprint(args[0])
		nodeMap, ok := args[1].(map[interface{}]interface{})

		// Check parameters

		if !ok {
			err = fmt.Errorf("Second parameter must be a map")
		}

		if err == nil && len(args) > 2 {
			if trans, ok = args[2].(graph.Trans); !ok {
				err = fmt.Errorf("Third parameter must be a transaction")
			}
		}

		// Build up node to store

		edge := data.NewGraphEdgeFromNode(NewGraphNodeFromECALMap(nodeMap))

		// Store the edge

		if err == nil {

			if trans != nil {
				err = trans.StoreEdge(part, edge)
			} else {
				err = f.GM.StoreEdge(part, edge)
			}
		}
	}

	return nil, err
}

/*
DocString returns a descriptive string.
*/
func (f *StoreEdgeFunc) DocString() (string, error) {
	return "Inserts or updates an edge in FishDB.", nil
}

/*
RemoveEdgeFunc removes an edge in FishDB.
*/
type RemoveEdgeFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *RemoveEdgeFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var err error

	if arglen := len(args); arglen != 3 && arglen != 4 {
		err = fmt.Errorf("Function requires 3 or 4 parameters: partition, edge key," +
			" edge kind and optionally a transaction")
	}

	if err == nil {
		var trans graph.Trans

		part := fmt.Sprint(args[0])
		key := fmt.Sprint(args[1])
		kind := fmt.Sprint(args[2])

		// Check parameters

		if len(args) > 3 {
			var ok bool

			if trans, ok = args[3].(graph.Trans); !ok {
				err = fmt.Errorf("Fourth parameter must be a transaction")
			}
		}

		// Remove the edge

		if err == nil {

			if trans != nil {
				err = trans.RemoveEdge(part, key, kind)
			} else {
				_, err = f.GM.RemoveEdge(part, key, kind)
			}
		}
	}

	return nil, err
}

/*
DocString returns a descriptive string.
*/
func (f *RemoveEdgeFunc) DocString() (string, error) {
	return "Removes an edge in FishDB.", nil
}

/*
FetchEdgeFunc fetches an edge in FishDB.
*/
type FetchEdgeFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *FetchEdgeFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var res interface{}
	var err error

	if arglen := len(args); arglen != 3 {
		err = fmt.Errorf("Function requires 3 parameters: partition, edge key and" +
			" edge kind")
	}

	if err == nil {
		var node data.Node

		part := fmt.Sprint(args[0])
		key := fmt.Sprint(args[1])
		kind := fmt.Sprint(args[2])

		conv := func(m map[string]interface{}) map[interface{}]interface{} {
			c := make(map[interface{}]interface{})
			for k, v := range m {
				c[k] = v
			}
			return c
		}

		// Fetch the node

		if node, err = f.GM.FetchEdge(part, key, kind); node != nil {
			res = conv(node.Data())
		}
	}

	return res, err
}

/*
DocString returns a descriptive string.
*/
func (f *FetchEdgeFunc) DocString() (string, error) {
	return "Fetches an edge in FishDB.", nil
}

/*
TraverseFunc traverses an edge in FishDB.
*/
type TraverseFunc struct {
	GM *graph.Manager
}

/*
Run executes the ECAL function.
*/
func (f *TraverseFunc) Run(instanceID string, vs parser.Scope, is map[string]interface{}, tid uint64, args []interface{}) (interface{}, error) {
	var res interface{}
	var err error

	if arglen := len(args); arglen != 4 {
		err = fmt.Errorf("Function requires 4 parameters: partition, node key," +
			" node kind and a traversal spec")
	}

	if err == nil {
		var nodes []data.Node
		var edges []data.Edge

		part := fmt.Sprint(args[0])
		key := fmt.Sprint(args[1])
		kind := fmt.Sprint(args[2])
		spec := fmt.Sprint(args[3])

		conv := func(m map[string]interface{}) map[interface{}]interface{} {
			c := make(map[interface{}]interface{})
			for k, v := range m {
				c[k] = v
			}
			return c
		}

		// Do the traversal

		if nodes, edges, err = f.GM.TraverseMulti(part, key, kind, spec, true); err == nil {

			resNodes := make([]interface{}, len(nodes))
			for i, n := range nodes {
				resNodes[i] = conv(n.Data())
			}
			resEdges := make([]interface{}, len(edges))
			for i, e := range edges {
				resEdges[i] = conv(e.Data())
			}
			res = []interface{}{resNodes, resEdges}
		}
	}

	return res, err
}

/*
DocString returns a descriptive string.
*/
func (f *TraverseFunc) DocString() (string, error) {
	return "Traverses an edge in FishDB from a given node.", nil
}
