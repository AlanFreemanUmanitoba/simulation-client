// display.admin.go
// handlers for actions specific to the admin

package display

import (
	"capfront/api"
	"capfront/models"
	"capfront/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Display the admin dashboard
func AdminDashboard(ctx *gin.Context) {
	username, err := synchWithServer(ctx)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf(" Could not retrieve user Data for user %s \n", username))
		return
	}

	ctx.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"Title": "Admin Dashboard",
		"users": models.AdminUserList,
	})
}

// Resets the main database
// Only available to admin
// TODO not done yet
func AdminReset(ctx *gin.Context) {
}

func Lock(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "choose-player.html", gin.H{
		"Title": "Choose player",
		"users": models.AdminUserList,
	})
}

func SelectUser(ctx *gin.Context) {
	u := ctx.Param("username")

	utils.Trace(utils.Yellow, fmt.Sprintf("user %s will play\n", u))
	// lock this user
	_, err := api.ServerRequest(models.Users[u].ApiKey, `admin/lock/`+u)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Could not play as user %s. It's just possible somebody else got in first", u))
		ctx.Abort()
		return
	}
	models.Users[u].IsLocked = true
	http.SetCookie(ctx.Writer, &http.Cookie{Name: "user", Value: u, Path: "/"})
	ctx.Request.URL.Path = `/`
	Router.HandleContext(ctx)
}
