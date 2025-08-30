/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package console

import (
	"fmt"
	"net/url"

	v1 "github.com/Fisch-Labs/FishDB/api/v1"
	"github.com/Fisch-Labs/common/stringutil"
)

// EQL Console
// ===========

/*
EQLConsole runs EQL queries.
*/
type EQLConsole struct {
	parent CommandConsoleAPI // Parent console API
}

/*
eqlConsoleKeywords are all keywords which this console can process.
*/
var eqlConsoleKeywords = []string{"part", "get", "lookup"}

/*
Run executes one or more commands. It returns an error if the command
had an unexpected result and a flag if the command was handled.
*/
func (c *EQLConsole) Run(cmd string) (bool, error) {

	if !cmdStartsWithKeyword(cmd, eqlConsoleKeywords) {
		return false, nil
	}

	// Escape query so it can be used in a request

	q := url.QueryEscape(cmd)

	resObj, err := c.parent.Req(
		fmt.Sprintf("%s%s?q=%s", v1.EndpointQuery, c.parent.Partition(), q), "GET", nil)

	if err == nil && resObj != nil {
		res := resObj.(map[string]interface{})
		var out []string

		header := res["header"].(map[string]interface{})

		labels := header["labels"].([]interface{})
		data := header["data"].([]interface{})
		rows := res["rows"].([]interface{})

		for _, l := range labels {
			out = append(out, fmt.Sprint(l))
		}
		for _, d := range data {
			out = append(out, fmt.Sprint(d))
		}
		for _, r := range rows {
			for _, c := range r.([]interface{}) {
				out = append(out, fmt.Sprint(c))
			}
		}

		c.parent.ExportBuffer().WriteString(stringutil.PrintCSVTable(out, len(labels)))
		fmt.Fprint(c.parent.Out(), stringutil.PrintGraphicStringTable(out, len(labels), 2, stringutil.SingleLineTable))
	}

	return true, err
}

/*
Commands returns an empty list. The command line is interpreted as an EQL query.
*/
func (c *EQLConsole) Commands() []Command {
	return nil
}
