package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Function to handle errors that need to be displayed to the user.
// Wrapped here to save space and also because error handling may
// be changed later.
func DisplayError(ctx *gin.Context, message string) {
	log.Output(1, message)
	ctx.HTML(http.StatusBadRequest, "errors.html", gin.H{
		"message": message,
	})
}

func DisplayLogin(ctx *gin.Context, message string) {
	log.Output(1, fmt.Sprint("Ask user to log in (possibly again) because of: ", message))
	ctx.HTML(http.StatusOK, "login.html", message)
}
