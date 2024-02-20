// display.actions.go
// handlers for actions requested by the user

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

// dispatches the action requested in the URL (of the form /actions/action)
func ActionHandler(c *gin.Context) {
	action := c.Param("action")

	url := api.APIsource + `action/` + action
	log.Output(1, fmt.Sprintf("Requesting remote server at URL %s to perform action %s\n", url, action))
	p := models.PostData{User: "admin"} //TODO get the actual user from the authorization package
	body, _ := json.Marshal(p)
	d := bytes.NewBuffer(body)
	r, err := http.NewRequest("GET", url, d)
	if err != nil {
		fmt.Println("!!!!!!!!There was an error composing the request. It was:", err)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Set("User-Agent", "Capitalism reader")
	r.Header.Add("Authorization", "Bearer "+auth.UserToken.Token)

	auth.VerifyTokenController(nil, r)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("!!!!!!!!!There was an error sending the request. It was:", err)
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Output(2, fmt.Sprintf("The response from the server was (%v)\n", string(b)))

	api.UserMessage = string(b)
	api.Refresh()
	ShowIndexPage(c)
	//TODO redisplay the page the user was looking at instead of just the index page
	//TODO SetTimeStamp(1)
}

// TODO remove boilerplate
func CreateSimulation(c *gin.Context) {
	template_id := c.Param("id")
	url := api.APIsource + `user/clone/` + template_id
	p := models.PostData{User: "alan"} //TODO get the actual user from the authorization package
	body, _ := json.Marshal(p)
	d := bytes.NewBuffer(body)

	log.Output(1, fmt.Sprintf("Requesting server to clone template %s using URL %s \n", template_id, url))

	r, err := http.NewRequest("GET", url, d)
	if err != nil {
		fmt.Println("!!!!!!!!There was an error composing the request. It was:", err)
		return
	}

	r.Header.Set("User-Agent", "Capitalism reader")
	r.Header.Add("Authorization", "Bearer "+auth.UserToken.Token)
	auth.VerifyTokenController(nil, r)

	r.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println("!!!!!!!!!There was an error sending the request. It was:", err)
		return
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	log.Output(2, fmt.Sprintf("The response from the server was (%v)\n", string(b)))

	api.UserMessage = string(b)
	api.Refresh()
	ShowIndexPage(c)

}
