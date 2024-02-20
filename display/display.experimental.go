// display.experimental.go
// handlers to display the objects of the simulation on the user's browser

package display

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DisplayMode(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "displaymode.html", gin.H{
		"Title": "Display Mode JS Tester",
	})

}
