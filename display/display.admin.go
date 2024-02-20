// display.admin.go
// handlers for actions specific to the admin

package display

import (
	"capfront/auth"
	"capfront/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Display the admin dashboard
func AdminDashboard(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/register")
		return
	}

	if username != "admin" {
		ctx.HTML(http.StatusOK, "errors.html", gin.H{
			"message": fmt.Errorf("only administrator can see the admin dashboard"),
		})
		return
	}
	ctx.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"Title":          "Admin Dashboard",
		"users":          models.Users,
		"username":       username,
		"loggedinstatus": loginStatus,
	})
}

// Resets the main database
// Only available to admin
func AdminReset(ctx *gin.Context) {
	username, _ := auth.Get_current_user(ctx)

	if username != "admin" {
		log.Output(1, fmt.Sprintf("User %s tried to reset the database", username))
		ShowIndexPage(ctx)
	}

	_, jsonErr := auth.ProtectedResourceServerRequest(username, "reset the database", `action/reset`)
	if jsonErr != nil {
		log.Output(1, "Reset failed")
	} else {
		log.Output(1, "COMPLETE RESET by admin")
	}

	AdminDashboard(ctx)
}
