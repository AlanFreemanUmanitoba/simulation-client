// display.actions.go
// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

package display

import (
	"capfront/api"
	"capfront/fetch"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// The state which follows each action.
var nextStates = map[string]string{
	`demand`:  `SUPPLY`,
	`supply`:  `TRADE`,
	`trade`:   `PRODUCE`,
	`produce`: `CONSUME`,
	`consume`: `INVEST`,
	`invest`:  `DEMAND`,
}

// pages for which redirection is OK.
func useLastVisited(last string) bool {
	switch last {
	case
		`/commodities`,
		`/industries`,
		`/classes`,
		`/stocks`,
		`/`:
		return true
	}
	return false
}

// Handles requests for the server to take an action comprising a stage
// of the circuit (demand,supply, trade, produce, invest), corresponding
// to a button press. This is specified by the URL parameter 'act'.
//
// Having requested the action from ths server, sets 'state' to the next
// stage of the circuit and redisplays whatever the user was looking at.
func ActionHandler(ctx *gin.Context) {
	// Comment for less detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Purple, fmt.Sprintf(" ActionHandler was called from %s #%d\n", file, no))
	}
	var param string
	err := ctx.ShouldBindUri(&param)
	if err != nil {
		fmt.Println("Malformed URL", err)
		ctx.String(http.StatusBadRequest, "Malformed URL")
		return
	}
	act := ctx.Param("action")
	username := utils.GUESTUSER
	userDatum := models.Users[username] // NOTE we assume the user exists in local storage
	lastVisitedPage := userDatum.LastVisitedPage
	log.Output(1, fmt.Sprintf("User %s wants the server to implement action %s. The last visited page was %s\n", username, act, lastVisitedPage))

	_, err = api.ServerRequest(userDatum.ApiKey, `action/`+act)
	if err != nil {
		utils.DisplayError(ctx, "The server could not complete the action")
	}

	// The action was taken. Now refresh the data from the server
	if !fetch.FetchUserObjects(ctx, username) {
		utils.DisplayError(ctx, "The server completed the action but did not send back any data.")
	}

	// Set the state so that the simulation can proceed to the next action.
	set_current_state(username, nextStates[act])

	// If the user was looking at a page that displays (but does not act),
	// redirect to it so the user can see the result of the action.
	// If not, redirect to the Index page.
	visitedPageURL := strings.Split(lastVisitedPage, "/")
	log.Output(1, fmt.Sprintf("The last page this user visited was %v and this was split into%v", lastVisitedPage, visitedPageURL))
	if useLastVisited(lastVisitedPage) {
		// utils.Trace(utils.Purple, fmt.Sprintf("User will be redirected to the last visited page which was %s\n", lastVisitedPage))
		ctx.Request.URL.Path = lastVisitedPage
	} else {
		// utils.Trace(utils.Purple, "User will be redirected to the Index Page, because the last visited URL was not a display page")
		ctx.Request.URL.Path = "/"
	}
	Router.HandleContext(ctx)
}

type CloneResult struct {
	Message       string `json:"message"`
	StatusCode    int    `json:"statusCode"`
	Simulation_id int    `json:"simulation_id"`
}

// Creates a new simulation for the user, from the template specified by the 'id' parameter.
// Initially, assume the user is 'guest'.
// This can be scaled up when and if login is introduced.
func CreateSimulation(ctx *gin.Context) {
	// Comment for shorter diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Green, fmt.Sprintf(" Clone Simulation was called from %s#%d\n", file, no))
	}

	username := utils.GUESTUSER
	user := models.Users[username] // Should we test for non-existing user?
	t := ctx.Param("id")
	id, _ := strconv.Atoi(t)
	log.Output(1, fmt.Sprintf("Creating a simulation from template %d for user %s", id, username))

	// Ask the server to create the clone and tell us the simulation id
	body, err := api.ServerRequest(user.ApiKey, `clone/`+t)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Failed to complete clone because of %v", err))
		return
	}

	// read the simulation id
	var result CloneResult
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Couldn't decode the clone result because of this error:%v", jsonErr))
		return
	}

	// Set the current simulation
	utils.Trace(utils.Green, fmt.Sprintf("Setting current simulation to be %d\n", result.Simulation_id))
	models.Users[username].CurrentSimulationID = result.Simulation_id

	// Diagnostic - comment or uncomment as needed
	// s, _ := json.MarshalIndent(models.Users[username], "  ", "  ")
	// fmt.Printf("User record after creating the simulation is %s\n", string(s))

	// Fetch the whole (new) dataset from the server
	// (until now we only told the server to create it - now we want it)
	if !fetch.FetchUserObjects(ctx, username) {
		utils.DisplayError(ctx, "WARNING: though the server created a simulation, we could not retrieve all its data")
	}
	// Initialise the timeStamp so that we are viewing the first dataset.
	// As the user moves through the circuit, this timestamp will move forwards.
	// Each time we move forward, a new dataset will be created.
	// This allows the user to view and compare with previous stages of the simulation.
	models.Users[username].ViewedTimeStamp = 0

	ctx.Request.URL.Path = "/"
	Router.HandleContext(ctx)
}
