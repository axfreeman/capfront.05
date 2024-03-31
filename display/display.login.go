// display.login.go
// Handles communication about authentication with the server (eg login, registration, etc)

package display

import (
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var genericAdvice = "Please contact the developer.\nPlease give as much information as you can, including the exact time this happened and the following message:\n"

// var rerouteMessage = gin.H{
// 	"message": "Could not log you in",
// 	"advice":  "Please try again",
// 	"info":    "Or ask me to register you",
// }

// List of possible error messages to display to the user when needed
var excuses = map[string]string{
	`client`:   `Sorry, there has been a programming error. ` + genericAdvice,
	`server`:   `Sorry, the server is down, or could not cope. ` + genericAdvice,
	`rejected`: `Sorry, the server did not accept this password. ` + genericAdvice,
	`narfy`:    `I'm afraid the server produced an incomprehensible response. ` + genericAdvice,
	`comms`:    `Sorry, I couldn't understand what the server said. ` + genericAdvice,
	`success`:  `Success`,
}

var userMessage string

// Displays a form to capture the user request to log in.
// The form specifies only one action, which is a submit button that POSTS the user name and password.
// This POST is handled by `ClientLoginRequest`.
func CaptureLoginRequest(ctx *gin.Context) {
	// uncomment for detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf(" The login form was called from %s#%d\n", file, no)
	}

	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"message": userMessage,
	})
}

// Service the form submitted when a user logs in.
// Because of the setup, it merely passes the request to the backend server
func HandleLoginRequest(ctx *gin.Context) {

	// extract the user name and the password that were submitted on the form.
	clientRequest := ctx.Request
	clientRequest.ParseForm()
	username := clientRequest.Form["username"][0]
	password := clientRequest.Form["password"][0]

	// ask the server to service this request and tell us the result
	token, result := ServerLogin(username, password)

	if token == nil {
		// Something went wrong; tell the user and offer some advice
		utils.DisplayError(ctx, fmt.Sprintf("%v", excuses[result]))
		return
	}

	// register the user name as a cookie in the user browser
	ctx.SetCookie("User", username, 34560000, "/", "", false, false)

	// pick up the local user record from the Users list
	userRecord, ok := models.Users[username]
	if !ok {
		// create a record if one does not exist
		log.Output(1, fmt.Sprintf("Creating a new user record for user %s", username))
		new_user := models.NewUserDatum(username)
		models.Users[username] = &new_user
		userRecord = models.Users[username]
	} else {
		log.Output(1, fmt.Sprintf("A user record already exists for user %s", username))
	}

	// Fill out user details from the server
	body, _ := auth.ProtectedResourceServerRequest(username, "Get user details ", `users/`+username)
	jsonErr := json.Unmarshal(body, &userRecord)

	if jsonErr != nil {
		// We couldn't understand the server's response
		utils.DisplayError(ctx, "We couldn't get your user details from the server")
		return
	}
	userRecord.Token = token.(string)

	if len(userRecord.History) != 0 {
		// Has the user already got a simulation going? If so, refresh it from the server.
		if !api.FetchUserObjects(ctx, username) {
			// In the unlikely case that the server authenticates this user but provides no data, tell the user.
			utils.DisplayError(ctx, "We couldn't get your data from the server")
			return
		}
	}

	// display the appropriate dashboard.
	if username == "admin" {
		ctx.Redirect(http.StatusMovedPermanently, "/admin/dashboard")
	} else {
		ctx.Redirect(http.StatusMovedPermanently, "/user/dashboard")
	}
}

// Compose and send a request to the server to log in.
//
// This function can be called either by the client (this project) using
// data entered by the user via the login form, or generated internally,
// to allow admin to retrieve information from the server.
//
//		username is the user supplied by the caller.
//		password is the password supplied by the caller.
//
//		returns a token if login is successful, nil otherwise.
//	  returns a string that summarises the result and indexes the 'excuses' list
func ServerLogin(username string, password string) (json.Token, string) {
	apiUrl := auth.APISOURCE + `auth/login`
	client := http.Client{Timeout: time.Second * 2}
	serverpayload := `username=` + username + `&password=` + password
	var req *http.Request
	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(serverpayload))
	if err != nil {
		return nil, "client"
	}
	req.Header.Set("Authorization", "Basic Og==")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, "server"
	}

	if res.StatusCode != 200 {
		return nil, "rejected"
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, "narfy"
	}

	var target map[string]string // Receives the token from the server
	jsonerr := json.Unmarshal(body, &target)
	if jsonerr != nil {
		return nil, "comms"
	}

	accessToken := target["access_token"]
	log.Output(1, fmt.Sprintf("The server has logged in user %s \n", username))
	return accessToken, fmt.Sprintf("Logged in user %s\n", username)
}

// logs the user out.
// error if the user is not logged in.
// This could arise if, for example, the token has expired and the user refreshes.
func ClientLogoutRequest(ctx *gin.Context) {
	username, err := auth.Get_current_user(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, fmt.Sprintf("Failed to log out because: %v", err))
		return
	}
	userDetails := models.Users[username]
	auth.ProtectedResourceServerRequest(username, "Log out", `auth/logout`)
	userDetails.Token = "invalid token"
	userDetails.LoggedIn = false
	CaptureLoginRequest(ctx)
}

// Asks the client to register.
func CaptureRegisterRequest(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}

// Service the form submitted when a user registers.
// Because of the setup, it merely passes the request to the backend server
func HandleRegisterRequest(ctx *gin.Context) {
	clientRequest := ctx.Request
	clientRequest.ParseForm()
	username := clientRequest.Form["username"][0]
	password := clientRequest.Form["password"][0]
	message, err := ServerRegister(username, password) // Ask the server to do the heavy lifting

	if err != nil { // something went wrong; tell the developer and tell the user
		utils.DisplayError(ctx, message)
		return
	}
	utils.DisplayLogin(ctx, "You can log in now ")

}

// Compose and send a request to the server to register.
func ServerRegister(username string, password string) (string, error) {
	apiUrl := auth.APISOURCE + `auth/register`
	client := http.Client{Timeout: time.Second * 2}
	serverpayload := `username=` + username + `&password=` + password
	serverRequest, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(serverpayload))
	if err != nil {
		return "client", err
	}
	serverRequest.Header.Set("Authorization", "Basic Og==")
	serverRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	print("Sending register request to server\n", serverpayload)
	res, err := client.Do(serverRequest)
	if err != nil {
		return "server", err
	}

	if res.StatusCode != 200 {
		return "rejected", err
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "narfy", readErr
	}

	var target models.ServerMessage // Receives the response from the server
	jsonerr := json.Unmarshal(body, &target)
	if jsonerr != nil {
		return "comms", jsonerr
	}

	if target.StatusCode != 200 {
		return "rejected", jsonerr
	}

	log.Output(1, fmt.Sprintf(" Registered user %s \n", username))

	new_user := models.NewUserDatum(username)
	models.Users[username] = &new_user
	return "Registration succeeded. Please log in", nil
}
