// api.server.go
// container for interacting with remote server

package api

import (
	"bytes"
	"capfront/colour"
	"capfront/logging"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var AccessToken string
var URLheader string

var Router *gin.Engine = gin.New()

// Runs once at startup.
// Retrieve users and templates from the server database.
func Initialise() {
	// Retrieve users on the server
	if !FetchAdminObject(utils.APISOURCE+`admin/users`, `users`) {
		log.Fatal("Could not retrieve user information from the server. Stopping")
	}
	for _, item := range models.AdminUserList {
		user := models.User{UserName: item.UserName, CurrentSimulationID: item.CurrentSimulationID, ApiKey: item.ApiKey}
		models.Users[item.UserName] = &user
	}

	// Retrieve the templates on the server
	if !FetchAdminObject(utils.APISOURCE+`templates/templates`, `templates`) {
		log.Fatal("Could not retrieve templates information from the server. Stopping")
	}
}

// Helper function to list out users and templates
func ListData() {
	fmt.Printf("\nTemplateList has %d elements which are:\n", len(models.TemplateList))
	for i := 0; i < len(models.TemplateList); i++ {
		fmt.Println(models.TemplateList[i])
	}
	m, _ := json.MarshalIndent(models.Users, " ", " ")
	fmt.Println(string(m))
}

// Prepare and send a request for a protected service to the server
// using the user's api key.
//
//	ctx is the context of a handler.
//	username is the name of the user requesting the service.
//	description is a user-friendly name for the requested action.
//	relativePath is appended to the server URL to tell the server what to do.
func ServerRequest(username string, description string, relativePath string) ([]byte, error) {
	_, file, no, ok := runtime.Caller(1)
	if ok {
		logging.Trace(colour.Cyan, fmt.Sprintf(" ServerRequest was called from %s#%d\n", file, no))
		logging.Trace(colour.Cyan, fmt.Sprintf(colour.Cyan+" Username was %s, description %s, relativePath %s, APIURL %s\n"+colour.Reset, username, description, relativePath, utils.APISOURCE))
	}

	user, ok := models.Users[username]
	if !ok {
		return nil, fmt.Errorf(" User %s is not in the local database", username)
	}

	url := utils.APISOURCE + relativePath
	body, _ := json.Marshal(models.RequestData{User: username}) // (overkill diagnostic? - not actually needed)
	resp, err := http.NewRequest("GET", url, bytes.NewBuffer(body))

	if err != nil {
		log.Output(1, fmt.Sprintf("Error %v for user %s from URL %s for resource %s \n", err, username, url, description))
		return nil, err
	}

	// logging.Trace(colour.Cyan, fmt.Sprintf(" Api key is [%s]\n", user.ApiKey))
	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Set("User-Agent", "Capitalism reader")
	resp.Header.Add("x-api-key", user.ApiKey)

	client := &http.Client{Timeout: time.Second * 5} // Timeout after 5 seconds
	res, _ := client.Do(resp)
	if res == nil {
		log.Output(1, "Server is down or misbehaving")
		return nil, nil
	}

	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		log.Output(1, fmt.Sprintf("Server rejected request '%s' with status %s\n", description, res.Status))
		logging.Trace(colour.Red, fmt.Sprintf("It said %s\n", string(b)))
		return nil, fmt.Errorf(string(b))
	}

	// Comment for fewer diagnostics
	fmt.Println(colour.Cyan + " Leaving ProtectedServerRequest. Everything seems to have worked." + colour.Reset)

	return b, nil
}
