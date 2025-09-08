/*
 * FishDB
 *
// Copyright 2025 Fisch-labs
 *
*/

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Fisch-Labs/FishDB/config"
	"github.com/Fisch-Labs/Toolkit/httputil"
)

const TESTPORT = ":9090"

var lastRes []string

type testEndpoint struct {
	*DefaultEndpointHandler
}

/*
handleSearchQuery handles a search query REST call.
*/
func (te *testEndpoint) HandleGET(w http.ResponseWriter, r *http.Request, resources []string) {
	lastRes = resources
	te.DefaultEndpointHandler.HandleGET(w, r, resources)
}

func (te *testEndpoint) SwaggerDefs(s map[string]interface{}) {
}

var testEndpointMap = map[string]RestEndpointInst{
	"/": func() RestEndpointHandler {
		return &testEndpoint{}
	},
}

func TestEndpointHandling(t *testing.T) {

	hs, wg := startServer()
	if hs == nil {
		return
	}
	defer func() {
		stopServer(hs, wg)
	}()

	queryURL := "http://localhost" + TESTPORT

	RegisterRestEndpoints(testEndpointMap)
	RegisterRestEndpoints(GeneralEndpointMap)

	// --- Test resource path parsing ---
	lastRes = nil
	if body, resp := sendTestRequestResponse(t, queryURL, "GET", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for GET: status %d, body %q", resp.StatusCode, body)
		return
	}
	if lastRes != nil {
		t.Error("Unexpected lastRes:", lastRes)
	}

	lastRes = nil
	if body, resp := sendTestRequestResponse(t, queryURL+"/foo/bar", "GET", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for /foo/bar: status %d, body %q", resp.StatusCode, body)
		return
	}
	if fmt.Sprint(lastRes) != "[foo bar]" {
		t.Error("Unexpected lastRes:", lastRes)
	}

	// Test trailing slashes
	lastRes = nil
	if body, resp := sendTestRequestResponse(t, queryURL+"/foo/bar/", "GET", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for /foo/bar/: status %d, body %q", resp.StatusCode, body)
		return
	}
	if fmt.Sprint(lastRes) != "[foo bar]" {
		t.Error("Unexpected lastRes:", lastRes)
	}

	// Test double slashes
	lastRes = nil
	if body, resp := sendTestRequestResponse(t, queryURL+"/foo//bar", "GET", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for /foo//bar: status %d, body %q", resp.StatusCode, body)
		return
	}
	if fmt.Sprint(lastRes) != "[foo bar]" {
		t.Errorf("Unexpected lastRes for double slash: %q", lastRes)
	}

	// --- Test unhandled HTTP methods ---
	if body, resp := sendTestRequestResponse(t, queryURL, "POST", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for POST: status %d, body %q", resp.StatusCode, body)
		return
	}
	if body, resp := sendTestRequestResponse(t, queryURL, "PUT", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for PUT: status %d, body %q", resp.StatusCode, body)
		return
	}
	if body, resp := sendTestRequestResponse(t, queryURL, "DELETE", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for DELETE: status %d, body %q", resp.StatusCode, body)
		return
	}
	if body, resp := sendTestRequestResponse(t, queryURL, "UPDATE", nil); body != "Method Not Allowed" || resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Unexpected response for UPDATE: status %d, body %q", resp.StatusCode, body)
		return
	}

	// --- Test /db/about endpoint ---
	body, resp := sendTestRequestResponse(t, queryURL+"/db/about", "GET", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK for /db/about, got %d", resp.StatusCode)
	}
	var aboutResponse struct {
		APIVersions []string `json:"api_versions"`
		Product     string   `json:"product"`
		Version     string   `json:"version"`
	}
	if err := json.Unmarshal([]byte(body), &aboutResponse); err != nil {
		t.Fatalf("Failed to decode /db/about JSON: %v\nbody: %s", err, body)
	}
	if aboutResponse.Product != "FishDB" {
		t.Errorf("Expected product 'FishDB', got %q", aboutResponse.Product)
	}
	if aboutResponse.Version != config.ProductVersion {
		t.Errorf("Expected version %q, got %q", config.ProductVersion, aboutResponse.Version)
	}
	if fmt.Sprint(aboutResponse.APIVersions) != "[v1]" {
		t.Errorf("Expected API versions '[v1]', got %v", aboutResponse.APIVersions)
	}

	// --- Test /db/swagger.json endpoint ---
	body, resp = sendTestRequestResponse(t, queryURL+"/db/swagger.json", "GET", nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK for /db/swagger.json, got %d", resp.StatusCode)
	}
	var swaggerDoc map[string]interface{}
	if err := json.Unmarshal([]byte(body), &swaggerDoc); err != nil {
		t.Fatalf("Failed to decode swagger.json: %v", err)
	}
	if swaggerDoc["swagger"] != "2.0" {
		t.Errorf("Expected swagger version '2.0', got %v", swaggerDoc["swagger"])
	}
	if swaggerDoc["basePath"] != "/db" {
		t.Errorf("Expected swagger basePath '/db', got %v", swaggerDoc["basePath"])
	}
	if swaggerDoc["host"] != APIHost {
		t.Errorf("Expected swagger host %q, got %v", APIHost, swaggerDoc["host"])
	}
	paths, ok := swaggerDoc["paths"].(map[string]interface{})
	if !ok || paths["/about"] == nil {
		t.Fatalf("Swagger doc is missing required '/about' path definition")
	}
}

/*
Send a request to a HTTP test server
*/
func sendTestRequest(t *testing.T, url string, method string, content []byte) string {
	body, _ := sendTestRequestResponse(t, url, method, content)
	return body
}

/*
Send a request to a HTTP test server
*/
func sendTestRequestResponse(t *testing.T, url string, method string, content []byte) (string, *http.Response) {
	var req *http.Request
	var err error

	if content != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(content))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	bodyStr := strings.Trim(string(body), " \n")

	// Try json decoding first

	out := bytes.Buffer{}
	err = json.Indent(&out, []byte(bodyStr), "", "  ")
	if err == nil {
		return out.String(), resp
	}

	// Just return the body

	return bodyStr, resp
}

/*
Start a HTTP test server.
*/
func startServer() (*httputil.HTTPServer, *sync.WaitGroup) {
	hs := &httputil.HTTPServer{}

	var wg sync.WaitGroup
	wg.Add(1)

	go hs.RunHTTPServer(TESTPORT, &wg)

	wg.Wait()

	// Server is started

	if hs.LastError != nil {
		panic(hs.LastError)
	}

	return hs, &wg
}

/*
Stop a started HTTP test server.
*/
func stopServer(hs *httputil.HTTPServer, wg *sync.WaitGroup) {
	if hs.Running {
		wg.Add(1)

		// Server is shut down

		hs.Shutdown()

		wg.Wait()

	} else {
		panic("Server was not running as expected")
	}
}
