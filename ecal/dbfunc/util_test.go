/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package dbfunc

import (
	"testing"

	"github.com/Fisch-Labs/FishDB/graph"
	"github.com/Fisch-Labs/ecal/interpreter"
	"github.com/Fisch-Labs/ecal/parser"
	"github.com/Fisch-Labs/ecal/util"
)

func TestRaiseGraphEventHandled(t *testing.T) {

	f := &RaiseGraphEventHandledFunc{}

	if _, err := f.DocString(); err != nil {
		t.Error(err)
		return
	}

	if _, err := f.Run("", nil, nil, 0, []interface{}{}); err != graph.ErrEventHandled {
		t.Error("Unexpected result:", err)
		return
	}
}

func TestRaiseWebEventHandled(t *testing.T) {

	f := &RaiseWebEventHandledFunc{}

	if _, err := f.DocString(); err != nil {
		t.Error(err)
		return
	}

	if _, err := f.Run("", nil, nil, 0, []interface{}{}); err == nil ||
		err.Error() != "Function requires 1 parameter: request response object" {
		t.Error(err)
		return
	}

	if _, err := f.Run("", nil, nil, 0, []interface{}{""}); err == nil ||
		err.Error() != "Request response object should be a map" {
		t.Error(err)
		return
	}

	astnode, _ := parser.ASTFromJSONObject(map[string]interface{}{
		"name": "foo",
	})

	_, err := f.Run("", nil, map[string]interface{}{
		"erp":     interpreter.NewECALRuntimeProvider("", nil, nil),
		"astnode": astnode,
	}, 0, []interface{}{map[interface{}]interface{}{}})

	if err.(*util.RuntimeErrorWithDetail).Type != ErrWebEventHandled {
		t.Error("Unexpected result:", err)
		return
	}
}
