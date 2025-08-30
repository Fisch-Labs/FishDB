/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

/*
Package api contains general REST API definitions.

The REST API provides an interface to FishDB. It allows querying and modifying
of the datastore. The API responds to GET, POST, PUT and DELETE requests in JSON
if the request was successful (Return code 200 OK) and plain text in all other cases.

# Common API definitions

/about

Endpoint which returns an object with version information.

	api_versions : List of available API versions e.g. [ "v1" ]
	product      : Name of the API provider (FishDB)
	version:     : Version of the API provider
	revision:    : Revision of the API provider

/swagger.json

Dynamically generated swagger definition file. See: http://swagger.io
*/
package api

import (
	"encoding/json"
	"net/http"

	"github.com/Fisch-Labs/FishDB/config"
)

/*
EndpointAbout is the about endpoint URL (rooted). Handles about/
*/
const EndpointAbout = APIRoot + "/about/"

/*
AboutEndpointInst creates a new endpoint handler.
*/
func AboutEndpointInst() RestEndpointHandler {
	return &aboutEndpoint{}
}

/*
Handler object for about operations.
*/
type aboutEndpoint struct {
	*DefaultEndpointHandler
}

/*
HandleGET returns about data for the REST API.
*/
func (a *aboutEndpoint) HandleGET(w http.ResponseWriter, r *http.Request, resources []string) {

	data := map[string]interface{}{
		"api_versions": []string{"v1"},
		"product":      "FishDB",
		"version":      config.ProductVersion,
	}

	// Write data

	w.Header().Set("content-type", "application/json; charset=utf-8")

	ret := json.NewEncoder(w)
	ret.Encode(data)
}
