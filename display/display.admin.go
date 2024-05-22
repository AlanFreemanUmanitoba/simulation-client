// display.admin.go
// handlers for actions specific to the admin

package display

import (
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
