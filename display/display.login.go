// display.login.go
// Handles communication about authentication with the server (eg login, registration, etc)

package display

import (
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var genericAdvice = "Please contact the developer.\nPlease give as much information as you can, including the exact time this happened and the following message:\n"

// container for building a user message
type apology struct {
	excuse string
}

// Generates an error report.
// `a` is a container for the explanation
// fault is the actual error
func (a apology) apologize(fault error) error {
	return fmt.Errorf(a.excuse+`%v`, fault)
}

// List of possible error messages to display to the user when needed
var excuses = map[string]apology{
	"client":   {"Sorry, there has been a programming error. " + genericAdvice},
	"server":   {"Sorry, the server is down, or could not cope. " + genericAdvice},
	"rejected": {"Sorry, the server is sulking. " + genericAdvice},
	"narfy":    {"I'm afraid the server produced an incomprehensible response. " + genericAdvice},
	"comms":    {"Sorry, I couldn't understand what the server said. " + genericAdvice},
}

var userMessage string

// Displays a form to capture the user request to log in.
// The form specifies only one action, which is a submit button that POSTS the user name and password.
// This POST is handled by `ClientLoginRequest`.
func CaptureLoginRequest(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.html", gin.H{
		"message": userMessage,
	})
}

// Service the form submitted when a user logs in.
// Because of the setup, it merely passes the request to the backend server
func HandleLoginRequest(ctx *gin.Context) {
	clientRequest := ctx.Request
	clientRequest.ParseForm()
	username := clientRequest.Form["username"][0]
	password := clientRequest.Form["password"][0]
	serverPayload, err := ServerLogin(username, password)

	if err != nil { // something went wrong; tell the developer and tell the user
		message := fmt.Sprintf("%s", serverPayload["message"])
		log.Output(1, message)
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"message": "Could not log you in",
			"advice":  "Please try again",
			"info":    "Or ask me to register you",
		})
		return
	}

	// register the user name as a cookie in the user browser
	// TODO fix up SameSite, wrong domain error, etc
	ctx.SetCookie("User", username, 34560000, "/", auth.APISOURCE, false, false)
	api.Refresh(ctx, username) // refresh the user's tables from the server at first login

	// Refresh user status from the server (which simulations we are using, etc)
	// TODO remove silly confusion between client URL 'user/' and server URL 'users/'
	body, _ := auth.ProtectedResourceServerRequest(username, " get user details ", `users/`+username)
	jsonErr := json.Unmarshal(body, &models.UserServerItem)

	if jsonErr != nil { // We couldn't understand the server's response
		// TODO display the error standardly as above and logout
		log.Output(1, "Failed to obtain user details for logged in user - cannot set current simulation right now")
	} else {
		log.Output(1, fmt.Sprintf("Setting current simulation to be %d", models.UserServerItem.CurrentSimulation))
		models.Users[username].CurrentSimulation = models.UserServerItem.CurrentSimulation
	}

	// display the appropriate dashboard.
	if username == "admin" {
		ctx.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		ctx.Redirect(http.StatusFound, "/user/dashboard")
	}
}

// Compose and send a request to the server to log in.
// This seems very laborious, more likely than not unnecessarily so.
// That's because it is a learning project for me.

// This function can be called either by the client (this project) using
// data entered by the user via the login form.
// OR can be generated internally, for example to log in to the server
// as admin and get some information from it.
func ServerLogin(username string, password string) (gin.H, error) {
	apiUrl := auth.APISOURCE + `auth/login`
	client := http.Client{Timeout: time.Second * 2}
	serverpayload := `username=` + username + `&password=` + password
	serverRequest, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(serverpayload))
	if err != nil {
		return map[string]any{"loggedinstatus": false, "message": excuses["client"].apologize(err)}, errors.New("login failed")
	}
	serverRequest.Header.Set("Authorization", "Basic Og==")
	serverRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(serverRequest)
	if err != nil {
		return map[string]any{"loggedinstatus": false, "message": excuses["server"].apologize(err)}, errors.New("login failed")
	}

	if res.StatusCode != 200 {
		return gin.H{"loggedinstatus": false, "message": excuses["rejected"].apologize(err)}, errors.New("login failed")
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return gin.H{"loggedinstatus": false, "message": excuses["narfy"].apologize(err)}, errors.New("login failed")
	}

	var target map[string]string // Receives the token from the server
	jsonerr := json.Unmarshal(body, &target)
	if jsonerr != nil {
		return gin.H{"loggedinstatus": false, "message": excuses["comms"].apologize(err)}, errors.New("login failed")
	}

	accessToken := target["access_token"]
	log.Output(1, fmt.Sprintf(" Logged in user %s \n", username))
	userDetails := models.Users[username]
	userDetails.Token = accessToken
	userDetails.LoggedIn = true // TODO think about cookie expiry and refresh
	return gin.H{"loggedinstatus": true, "message": fmt.Sprintf("Logged in user %s\n", username)}, nil
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
	userDetails.LoggedIn = false // TODO think about cookie expiry and refresh
	CaptureLoginRequest(ctx)
}

// Asks the client to register.
// TODO integrate this with the login page
func CaptureRegisterRequest(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "register.html", gin.H{})
}

// Service the form submitted when a user logs in.
// Because of the setup, it merely passes the request to the backend server
func HandleRegisterRequest(ctx *gin.Context) {
	clientRequest := ctx.Request
	clientRequest.ParseForm()
	username := clientRequest.Form["username"][0]
	password := clientRequest.Form["password"][0]
	serverPayload, err := ServerRegister(username, password)

	if err != nil { // something went wrong; tell the developer and tell the user
		message := fmt.Sprintf("%s", serverPayload["message"])
		log.Output(1, message)
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"message": "Could not log you in",
			"advice":  "Please try again",
			"info":    "Or ask me to register you",
		})
		return
	}
}

// Compose and send a request to the server to register.
// This seems very laborious, more likely than not unnecessarily so.
// It also has a lot of boilerplate code that simply repeats what is
// in CaptureLoginRequest().
// That's because it is a learning project for me.

// This function can be called either by the client (this project) using
// data entered by the user via the login form.
// OR can be generated internally, though I can't think of a user case for that.
func ServerRegister(username string, password string) (gin.H, error) {
	apiUrl := auth.APISOURCE + `auth/register`
	client := http.Client{Timeout: time.Second * 2}
	serverpayload := `username=` + username + `&password=` + password
	serverRequest, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(serverpayload))
	if err != nil {
		return map[string]any{"loggedinstatus": false, "message": excuses["client"].apologize(err)}, errors.New("registration failed")
	}
	serverRequest.Header.Set("Authorization", "Basic Og==")
	serverRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(serverRequest)
	if err != nil {
		return map[string]any{"loggedinstatus": false, "message": excuses["server"].apologize(err)}, errors.New("registration failed")
	}

	if res.StatusCode != 200 {
		return gin.H{"loggedinstatus": false, "message": excuses["rejected"].apologize(err)}, errors.New("registration failed")
	}
	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return gin.H{"loggedinstatus": false, "message": excuses["narfy"].apologize(err)}, errors.New("registration failed")
	}

	var target map[string]string // Receives the token from the server
	jsonerr := json.Unmarshal(body, &target)
	if jsonerr != nil {
		return gin.H{"loggedinstatus": false, "message": excuses["comms"].apologize(err)}, errors.New("registration failed")
	}

	log.Output(1, fmt.Sprintf(" Registered user %s \n", username))
	return gin.H{"message": "Registration succeeded. Please log in"}, errors.New("registration failed")
}
