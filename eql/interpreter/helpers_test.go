/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package interpreter

import (
	"testing"

	"github.com/Fisch-Labs/FishDB/eql/parser"
)

func TestHelperRuntime(t *testing.T) {
	gm, _ := simpleGraph()
	rt := NewGetRuntimeProvider("test", "main", gm, &testNodeInfo{&defaultNodeInfo{gm}})

	// Test simple value runtime

	ast, err := parser.ParseWithRuntime("test", "get mynode", rt)
	if err != nil {
		t.Error(err)
		return
	}

	if val, _ := ast.Children[0].Runtime.Eval(); val != "mynode" {
		t.Error("Unexpected eval result:", val)
		return
	}

	if err := ast.Children[0].Runtime.Validate(); err != err {
		t.Error(err)
		return
	}

	// Test not implemented runtime

	irt := invalidRuntimeInst(rt.eqlRuntimeProvider, ast.Children[0])

	if err := irt.Validate(); err.Error() != "EQL error in test: Invalid construct (value) (Line:1 Pos:5)" {
		t.Error("Unexpected validate result:", err)
		return
	}

	if _, err := irt.Eval(); err.Error() != "EQL error in test: Invalid construct (value) (Line:1 Pos:5)" {
		t.Error("Unexpected validate result:", err)
		return
	}
}
