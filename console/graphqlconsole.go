/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package console

import (
	"encoding/json"
	"fmt"

	v1 "github.com/Fisch-Labs/FishDB/api/v1"
	"github.com/Fisch-Labs/common/errorutil"
)

// GraphQL Console
// ===============

/*
GraphQLConsole runs GraphQL queries.
*/
type GraphQLConsole struct {
	parent CommandConsoleAPI // Parent console API
}

/*
graphQLConsoleKeywords are all keywords which this console can process.
*/
var graphQLConsoleKeywords = []string{"{", "query", "mutation"}

/*
Run executes one or more commands. It returns an error if the command
had an unexpected result and a flag if the command was handled.
*/
func (c *GraphQLConsole) Run(cmd string) (bool, error) {

	if !cmdStartsWithKeyword(cmd, graphQLConsoleKeywords) {
		return false, nil
	}

	q, err := json.Marshal(map[string]interface{}{
		"operationName": nil,
		"variables":     nil,
		"query":         cmd,
	})
	errorutil.AssertOk(err)

	resObj, err := c.parent.Req(
		fmt.Sprintf("%s%s", v1.EndpointGraphQL, c.parent.Partition()), "POST", q)

	if err == nil && resObj != nil {

		actualResultBytes, _ := json.MarshalIndent(resObj, "", "  ")
		out := string(actualResultBytes)

		c.parent.ExportBuffer().WriteString(out)
		fmt.Fprint(c.parent.Out(), out)
	}

	return true, err
}

/*
Commands returns an empty list. The command line is interpreted as a GraphQL query.
*/
func (c *GraphQLConsole) Commands() []Command {
	return nil
}
