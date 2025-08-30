/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package util

import (
	"errors"
	"testing"
)

func TestGraphError(t *testing.T) {
	err := GraphError{errors.New("TestError"), ""}

	if err.Error() != "GraphError: TestError" {
		t.Error("Unexpected result", err.Error())
		return
	}

	err = GraphError{errors.New("TestError"), "SomeDetail"}

	if err.Error() != "GraphError: TestError (SomeDetail)" {
		t.Error("Unexpected result", err.Error())
		return
	}
}
