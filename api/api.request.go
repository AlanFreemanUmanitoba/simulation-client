// api.server.go
// container for interacting with remote server

package api

import (
	"bytes"
	"capfront/utils"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Prepare and send a request for a protected service to the server
// using the user's api key.
//
//		apiKey is the key
//		url is appended to apiSource to tell the server what to do.
//
//	 Returns: byte array with the server response
//	 Returns: error if anything went wrong, or nil
func ServerRequest(apiKey string, url string) ([]byte, error) {
	utils.Trace(utils.Cyan, fmt.Sprintf("Entering ServerRequest with apiKey %s and relative path %s\n", apiKey, url))
	resp, err := http.NewRequest("GET", utils.APISOURCE+url, bytes.NewBuffer([]byte(`{"origin":"Simulation-client"}`)))
	if err != nil {
		utils.Trace(utils.Red, "Malformed client request")
		return nil, err
	}

	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Set("User-Agent", "Capitalism reader")
	resp.Header.Add("x-api-key", apiKey)

	client := &http.Client{Timeout: time.Second * 5} // Timeout after 5 seconds
	res, _ := client.Do(resp)
	if res == nil {
		utils.Trace(utils.Red, "Server is down or misbehaving\n")
		return nil, errors.New("server did not respond")
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		utils.Trace(utils.Red, fmt.Sprintf("Server rejected the request with status %s\n", res.Status))
		utils.Trace(utils.Red, fmt.Sprintf("It said %s\n", string(b)))
		return nil, errors.New(string(b))
	}
	// utils.Trace(utils.Cyan, "Leaving ServerRequest, everything looks good so far\n")
	return b, nil
}
