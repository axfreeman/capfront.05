// auth.common.go
// container for variables and types needed for authentication

package auth

import (
	"bytes"
	"capfront/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var AccessToken string
var URLheader string

// TODO put these into an env file
var SECRET_ADMIN_PASSWORD string = "insecure"

var APISOURCE = `https://www.datapaedia.org/` // Comment for production version
// var APISOURCE = `http://127.0.0.1:8000/` // Alternate for local version

// Get the user's identity from the browser by retrieving a cookie
// return err if cookie cannot be retrieved (normally, because it isn't there)
func Get_current_user(ctx *gin.Context) (string, error) {
	userCookie, err := ctx.Request.Cookie("User")
	if err != nil {
		return "", err
	}

	// fmt.Println("Cookie returned:", userCookie.Value) // Uncomment for verbose diagnostics
	username := userCookie.Value
	return username, nil
}

// Helper function to prepare and send a request for a protected service to the server
// ctx is the context of a handler
// username is the name of the user requesting the service
// description is a user-friendly name for the action being requested, which is used to produce error messages
// relativePath is appended to the URL of the remote server and tells the server what we want it to do
func ProtectedResourceServerRequest(username string, description string, relativePath string) ([]byte, error) {
	user, ok := models.Users[username]

	if !ok {
		log.Output(1, fmt.Sprintf("Attempt to access the server by non-existent user %s", username))
		return nil, fmt.Errorf("user %s tried to access the server, but we don't have any record of that user", username)
	}

	accessToken := user.Token
	userMessage := models.UserMessage{StatusCode: http.StatusTeapot, Message: ""} // Couldn't resist it
	user.UserMessage = &userMessage
	url := APISOURCE + relativePath
	log.Output(1, fmt.Sprintf("User %s asked URL %s for resource %s \n", username, relativePath, description))

	body, _ := json.Marshal(models.RequestData{User: username}) // Wrap username in RequestData struct to prepare for unmarshal
	resp, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		log.Output(1, fmt.Sprintf("Error %v for user %s from URL %s for resource %s \n", err, username, url, description))
		userMessage.StatusCode = http.StatusBadRequest
		userMessage.Message = fmt.Sprintf("Server error %v ", err)
		return nil, err
	}

	resp.Header.Add("Content-Type", "application/json")
	resp.Header.Set("User-Agent", "Capitalism reader")
	resp.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: time.Second * 2} // Timeout after 2 seconds
	res, _ := client.Do(resp)
	if res == nil {
		// Server failure
		// TODO display nice error screen
		log.Output(1, "Server is down or misbehaving")
		return nil, nil
	}

	if res.StatusCode != 200 {
		log.Output(1, fmt.Sprintf("Server at URL %s rejected user %s's request '%s'. Status code was %s\n", relativePath, username, description, res.Status))
		userMessage.StatusCode = http.StatusBadRequest
		userMessage.Message = fmt.Sprintf("Server error %v ", err)
		return nil, fmt.Errorf("could not access resource %s", description)
	}

	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)

	userMessage.StatusCode = http.StatusOK
	// The content of the user message, if the action succeeds,
	// should be set by the handler responsible for it.
	return b, nil
}
