// display.objects.go
// handlers to display the objects of the simulation on the user's browser
// TODO move this to 'display.objects.go'

package display

import (
	"capfront/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// display all industry stocks in the current simulation
func ShowIndustryStocks(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "industry_stocks.html", gin.H{
		"Title":          "Industry Stocks",
		"stocks":         models.Users[username].IndustryStockList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all the class stocks in the current simulation
func ShowClassStocks(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "class_stocks.html", gin.H{
		"Title":          "Class Stocks",
		"stocks":         models.Users[username].ClassStockList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}
