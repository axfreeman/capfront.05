// display.login.go
// Handles communication about authentication with the server (eg login, registration, etc)

package display

import (
	"bytes"
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// temporary fix so that whenever the frontend starts up, it logs in as a known user
// for testing purposes - saves having to do it each time we run it
func Fake_login() {
	username := `alan`
	password := `insecure`
	fmt.Println(username, password)
	DoLogin(username, password)
}

// given a username and a password, tell the server to log me in
// TODO shorten this - here it is done at length so I can see what is going on
func DoLogin(username string, password string) {
	userDetails := models.UserLoginDetails{UserName: username, Password: password}

	log.Output(1, fmt.Sprintf("Client Request to log in from user %s ", userDetails.UserName))
	apiRequest := api.APIsource + `auth/login`
	log.Output(1, fmt.Sprintf("API call to post to URL %s ", apiRequest))

	jsonStr, _ := json.Marshal(userDetails)
	// log.Output(1, fmt.Sprintf("json string is %v ", jsonStr))

	req, err := http.NewRequest("POST", apiRequest, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Output(1, fmt.Sprintf("Failed to create a valid login request to remote server because: %v", err))
	}
	req.Header.Set("X-Custom-Header", "Login request")
	req.Header.Set("Content-Type", "application/json")

	// log.Output(1, fmt.Sprintf("About to send a login request that looks like this: %v", req))
	// log.Output(1, fmt.Sprintf("the method is: %v", req.Method))
	log.Output(1, fmt.Sprintf("Login request sent: %v", req.Body))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Output(1, fmt.Sprintf("API Post to remote server failed because of %v", err))
	}
	log.Output(1, "API Post was accepted - now being tested for validity")
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Th response Body was:", string(body))

	if err := json.Unmarshal(body, &auth.UserToken); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Cannot unmarshal the JSON which the server sent back")
	}

	// fmt.Println("the token was", auth.UserToken.Token)
	log.Output(1, fmt.Sprintf("User %s is now logged in ", username))
	auth.LoggedInUser = username

	// The above integrates the details specified in the below; retained for now for reference
	// 	curl -X 'POST' \
	//   'http://localhost:8000/auth/login' \
	//   -H 'accept: application/json' \
	//   -H 'Content-Type: application/json' \
	//   -d '{
	//   "username": "string",
	//   "password": "string"
	// }'

}
func ServiceLoginRequest(ctx *gin.Context) {
	r := ctx.Request
	// diagnostics
	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]
	DoLogin(username, password)
	api.Refresh()
	ShowIndexPage(ctx)
}

// displays a form to capture the user request
// the form specifies only one action, which is a submit button that POSTS the user name and password
// this POST is handled by `LoginAtServer`
func CaptureLoginRequest(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "forms.html", gin.H{})
}
