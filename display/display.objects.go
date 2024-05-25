// display.objects.go
// handlers to display the objects of the simulation on the user's browser

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

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine = gin.New()

// Authentication. Should eventually become middleware.
//
// Maintains synchronisation between data on the server and on the client.
//
//	If the client does not know enough to establish a connection, redirects
//	to login and aborts, so that nothing else happens for the rest of the
//	handler chain.
//
//	If the server does not respond or returns an error, redirects to an
//	error page and aborts.
//
//	If the server responds correctly, returns the username and any error.
//
// Sets 'LastVisitedPage' so we can return here after an action.
func synchWithServer(ctx *gin.Context) (string, error) {
	// Comment for less detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Yellow, fmt.Sprintf(" SynchWithServer was called from %s#%d\n", file, no))
	}

	// Check whether the server is responding at all.
	// TODO some code duplication here.
	body, err := api.ServerRequest(models.Users["admin"].ApiKey, `admin/user/admin`)
	if (err != nil) || (body == nil) {
		errmsg := fmt.Sprintf("Sorry, the server is not functioning :%v\n", err)
		utils.Trace(utils.Red, errmsg)
		utils.DisplayError(ctx, errmsg)
		ctx.Abort()
		return "admin", err
	}

	// find out what the browser knows
	userCookie, err := ctx.Request.Cookie("user")
	if err != nil {
		utils.Trace(utils.BrightMagenta, fmt.Sprintf("Could not retrieve any user cookie becase of %v\n", err))
		ctx.Request.URL.Path = "/admin/choose-players"
		Router.HandleContext(ctx)
		ctx.Abort()
		return "Unknown User", err
	}

	// Comment for briefer diagnostics
	fmt.Println("Cookie returned:", userCookie.Value)
	// TODO username := userCookie.Value
	username := userCookie.Value
	user := models.Users[username]
	//TODO this ought to be impossible, so maybe this test is overkill
	if user.ApiKey == "" {
		utils.Trace(utils.BrightMagenta, "ERROR: User has no api key\n")
		ctx.Request.URL.Path = "/admin/choose-players"
		Router.HandleContext(ctx)
		ctx.Abort()
		return "No API key", err
	}

	// find out what the server knows
	//TODO this ought to be impossible, so maybe this test is overkill
	body, err = api.ServerRequest(user.ApiKey, `admin/user/`+username)
	if err != nil {
		utils.Trace(utils.BrightMagenta, fmt.Sprintf("The server was not functioning :%v\n", err))
		utils.DisplayError(ctx, "The API server is not functioning")
		ctx.Abort()
		return username, err
	}
	synched_user := new(models.User)          // Isolated user record, not added to the list of users.
	err = json.Unmarshal(body, &synched_user) // it's only there as a receptacle for the server data.
	if err != nil {
		utils.Trace(
			utils.BrightMagenta,
			fmt.Sprintf("Couldn't make sense of what the server said about user %s :%v\n", username, err))
		ctx.Request.URL.Path = "/admin/choose-players"
		Router.HandleContext(ctx)
		ctx.Abort()
		return username, err
	}
	userDetails, _ := json.MarshalIndent(synched_user, " ", " ")
	utils.Trace(utils.BrightMagenta, fmt.Sprintf("The server sent this user record: %s\n", string(userDetails)))

	utils.Trace(utils.BrightMagenta, fmt.Sprintf(
		"The server sent a user record which says the current simulation is %d; the client says it is %d\n",
		synched_user.CurrentSimulationID,
		models.Users[username].CurrentSimulationID,
	))

	// Update client data from the server if need be.
	if models.Users[username].CurrentSimulationID != synched_user.CurrentSimulationID {
		utils.Trace(utils.Yellow, fmt.Sprintf("We are out of synch. Server thinks our simulation is %d and client says it is %d\n",
			synched_user.CurrentSimulationID,
			models.Users[username].CurrentSimulationID))
		// Resynchronise
		if !fetch.FetchUserObjects(ctx, username) {
			utils.Trace(utils.Red, fmt.Sprintf("ERROR: Could not retrieve data for user %s\n", username))
			ctx.Request.URL.Path = "/admin/choose-players"
			Router.HandleContext(ctx)
			ctx.Abort()
			return username, fmt.Errorf("could not retrieve data for user %s", username)
		}
	}

	models.Users[username].LastVisitedPage = ctx.Request.URL.Path
	utils.Trace(utils.BrightMagenta, fmt.Sprintf("User %s is good to go\n", username))
	return username, nil
}

// Helper function to list out users and templates
func ListData() {
	fmt.Printf("\nTemplateList has %d elements which are:\n", len(models.TemplateList))
	for i := 0; i < len(models.TemplateList); i++ {
		fmt.Println(models.TemplateList[i])
	}

	fmt.Printf("\nAdminUserList has %d elements which are:\n", len(models.AdminUserList))
	for i := 0; i < len(models.AdminUserList); i++ {
		fmt.Println(models.AdminUserList[i])
	}

	fmt.Println("\nUsers", len(models.Users))
	m, _ := json.MarshalIndent(models.Users, " ", " ")
	fmt.Println(string(m))

}

// helper function to obtain the state of the current simulation
// to be replaced by inline call
func Get_current_state(username string) string {
	return models.Users[username].Get_current_state()
}

// helper function to set the state of the current simulation.
// To be replaced by inline call
func set_current_state(username string, new_state string) {
	u, ok := models.Users[username]
	if !ok {
		return
	}
	u.Set_current_state(new_state)
}

// display all commodities in the current simulation
func ShowCommodities(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Commodities")
		return
	}
	user := models.Users[username]
	state := user.Get_current_state()
	commodityViews := user.CommodityViews()

	ctx.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":          "Commodities",
		"commodities":    user.Commodities(),
		"commodityViews": commodityViews,
		"username":       username,
		"state":          state,
	})
}

// display all industries in the current simulation
func ShowIndustries(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Industries ")
		return
	}

	user := models.Users[username]
	state := user.Get_current_state()
	industryViews := user.IndustryViews()

	ctx.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":         "Industries",
		"industries":    user.Industries(),
		"industryViews": industryViews,
		"username":      username,
		"state":         state,
	})
}

// display all classes in the current simulation
func ShowClasses(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Classes ")
		return
	}

	user := models.Users[username]
	state := models.Users[username].Get_current_state()
	classViews := user.ClassViews()

	classViewAsString, _ := json.MarshalIndent(classViews, " ", " ")
	utils.Trace(utils.Cyan, fmt.Sprintf("Class Views:\n%s\n ", string(classViewAsString)))

	ctx.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":      "Classes",
		"classes":    models.Users[username].Classes(),
		"classViews": classViews,
		"username":   username,
		"state":      state,
	})
}

// Display one specific commodity
func ShowCommodity(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one Commodity ")
		return
	}

	state := models.Users[username].Get_current_state()
	id, _ := strconv.Atoi(ctx.Param("id"))

	clist := *models.Users[username].Commodities()
	for i := 0; i < len(clist); i++ {
		if id == clist[i].Id {
			ctx.HTML(http.StatusOK, "commodity.html", gin.H{
				"Title":     "Commodity",
				"commodity": clist[i],
				"username":  username,
				"state":     state,
			})
		}
	}
}

// Display one specific industry
func ShowIndustry(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one industry ")
		return
	}

	state := models.Users[username].Get_current_state()
	id, _ := strconv.Atoi(ctx.Param("id"))
	ilist := *models.Users[username].Industries()
	for i := 0; i < len(ilist); i++ {
		if id == ilist[i].Id {
			ctx.HTML(http.StatusOK, "industry.html", gin.H{
				"Title":    "Industry",
				"industry": ilist[i],
				"username": username,
				"state":    state,
			})
		}
	}
}

// Display one specific class
func ShowClass(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display one class ")
		return
	}

	state := models.Users[username].Get_current_state()
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	list := *models.Users[username].Classes()

	for i := 0; i < len(list); i++ {
		if id == list[i].Id {
			ctx.HTML(http.StatusOK, "class.html", gin.H{
				"Title":    "Class",
				"class":    list[i],
				"username": username,
				"state":    state,
			})
		}
	}
}

// Displays snapshot of the economy

func ShowIndexPage(ctx *gin.Context) {
	// Uncomment for more detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		utils.Trace(utils.Yellow, fmt.Sprintf(" ShowIndexPage was called from %s#%d\n", file, no))
	}
	username, err := synchWithServer(ctx)
	if err != nil {
		ctx.Abort()
		return
	}

	u := models.Users[username]
	state := u.Get_current_state()
	clist := u.Commodities()
	ilist := u.Industries()
	cllist := u.Classes()
	commodityViews := u.CommodityViews()
	industryViews := u.IndustryViews()
	classViews := u.ClassViews()

	// industryViewAsString, _ := json.MarshalIndent(industryViews, " ", " ")
	// utils.Trace(utils.BrightCyan, "  Industry view before displaying index page is\n"+string(industryViewAsString)+"/n")

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title":          "Economy",
		"industries":     ilist,
		"commodities":    clist,
		"commodityViews": commodityViews,
		"industryViews":  industryViews,
		"classes":        cllist,
		"classViews":     classViews,
		"username":       username,
		"state":          state,
	})
}

// Fetch the trace from the local database
func ShowTrace(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Trace records ")
		return
	}

	state := models.Users[username].Get_current_state()
	tlist := *models.Users[username].Traces(models.Users[username].ViewedTimeStamp)

	ctx.HTML(
		http.StatusOK,
		"trace.html",
		gin.H{
			"Title":    "Simulation Trace",
			"trace":    tlist,
			"username": username,
			"state":    state,
		},
	)
}

// Display all templates, and all simulations belonging to this user,
// in the user dashboard.
func UserDashboard(ctx *gin.Context) {
	if _, file, no, ok := runtime.Caller(1); ok {
		utils.Trace(utils.Yellow, fmt.Sprintf(" User Dashboard was called from %s line #%d\n", file, no))
	}

	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display the User dashboard ")
		return
	}

	state := models.Users[username].Get_current_state()
	slist := *models.Users[username].Simulations()

	ctx.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":       "Dashboard",
		"simulations": slist,
		"templates":   models.TemplateList,
		"username":    username,
		"state":       state,
	})
}

// Diagnostic endpoint to display the data in the system.
func DataHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.Users)
}

// TODO not working yet
func SwitchSimulation(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Data Listing ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to switch to simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

// TODO not working yet
func DeleteSimulation(ctx *gin.Context) {
}

// TODO not working yet
func RestartSimulation(ctx *gin.Context) {
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display industry Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to show industry stocks %d", username, id))

	state := models.Users[username].Get_current_state()
	islist := *models.Users[username].IndustryStocks(models.Users[username].ViewedTimeStamp)

	ctx.HTML(http.StatusOK, "industry_stocks.html", gin.H{
		"Title":    "Industry Stocks",
		"stocks":   islist,
		"username": username,
		"state":    state,
	})
}

// display all the class stocks in the current simulation
func ShowClassStocks(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, " Could not retrieve data from Server while trying to display Class Stocks ")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to show class stocks %d", username, id))
	state := models.Users[username].Get_current_state()
	cslist := *models.Users[username].ClassStocks(models.Users[username].ViewedTimeStamp)

	ctx.HTML(http.StatusOK, "class_stocks.html", gin.H{
		"Title":    "Class Stocks",
		"stocks":   cslist,
		"username": username,
		"state":    state,
	})
}
