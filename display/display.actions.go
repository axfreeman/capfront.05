// display.actions.go
// handlers for actions requested by the user

package display

// This module processes the actions that take the simulation through
// a circuit - Demand, Supply, Trade, Produce, Consume, Invest

import (
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"capfront/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Dispatches the action requested in the URL
// Eg /actions/trade will perform the trade action, and so on.
type Action struct {
	A string `uri:"action"`
}

// Handles requests for the server to take an action comprising a stage
// of the circuit (demand,supply, trade, produce, invest), corresponding
// to a button press. This is specified by the URL parameter 'act'
// Having requested the action from ths server, sets 'state' to the next
// stage of the circuit and redisplays whatever the user was looking at
func ActionHandler(ctx *gin.Context) {
	// Uncomment for more detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf(" ActionHandler was called from %s #%d\n", file, no)
	}
	var param Action
	err := ctx.ShouldBindUri(&param)
	if err != nil {
		fmt.Println("Malformed URL", err)
		ctx.String(http.StatusBadRequest, "Malformed URL")
		return
	}
	act := ctx.Param("action")
	username, _ := auth.Get_current_user(ctx)
	userDatum, ok := models.Users[username]
	if !ok {
		// This can happen if, for example, the user goes for a coffee...
		utils.DisplayError(ctx, "The server data has changed since you last visited. Please log in again")
	}

	lastVisitedPage := userDatum.LastVisitedPage
	log.Output(1, fmt.Sprintf("User %s wants the server to implement action %s. The last visited page was %s\n", username, act, lastVisitedPage))

	_, err = auth.ProtectedResourceServerRequest(username, act, `action/`+act)
	if err != nil {
		utils.DisplayError(ctx, "The server could not complete the action")
	}

	// The action was taken. Now refresh the data from the server
	// TODO move the timestamp forward and create a new history item.
	if !api.FetchUserObjects(ctx, username) {
		utils.DisplayError(ctx, "The server completed the action but did not send back any data.")
	}

	// TODO use the state information supplied by the server - this code duplicates the server's prerogative
	user := models.Users[username]
	switch act {
	case "demand":
		set_current_state(username, "SUPPLY")
		user.UserMessage.Message = "Demand Complete - watch this space"
	case "supply":
		set_current_state(username, "TRADE")
		user.UserMessage.Message = "Supply Complete - watch this space"
	case "trade":
		set_current_state(username, "PRODUCE")
		user.UserMessage.Message = "Trade Complete - watch this space"
	case "produce":
		set_current_state(username, "CONSUME")
		user.UserMessage.Message = "Production Complete - watch this space"
	case "consume":
		set_current_state(username, "INVEST")
		user.UserMessage.Message = "Consumption Complete - watch this space"
	case "invest":
		set_current_state(username, "DEMAND")
		user.UserMessage.Message = "Investment is not yet coded"
	default:
		set_current_state(username, "UNKNOWN")
		user.UserMessage.Message = "There has been a programme error of some kind"
	}
	// If the user has just visited a page that displays (but does not act!!!!), redirect to it.
	// If not, redirect to the Index page
	// This is a very crude mechanism
	visitedPageURL := strings.Split(lastVisitedPage, "/")
	log.Output(1, fmt.Sprintf("The last page this user visited was %v and this was split into%v", lastVisitedPage, visitedPageURL))
	// v := visitedPageURL[0]
	if lastVisitedPage == `/commodities` || lastVisitedPage == `/industries` || lastVisitedPage == `/classes` || lastVisitedPage == `/stocks` {
		log.Output(1, fmt.Sprintf("User will be redirected to the last visited page which was %s\n", lastVisitedPage))
		ctx.Redirect(http.StatusMovedPermanently, lastVisitedPage)
	} else {
		log.Output(1, "The user will be redirected to the Index Page, because the last visited URL was not a display page")
		ctx.Redirect(http.StatusMovedPermanently, "/index")
	}
}

type CloneResult struct {
	Message       string `json:"message"`
	Simulation_id int    `json:"simulation"`
}

// Creates a new simulation for the logged-in user, from the template specified by the 'id' parameter
func CreateSimulation(ctx *gin.Context) {
	username, _ := auth.Get_current_user(ctx)
	t := ctx.Param("id")
	id, _ := strconv.Atoi(t)
	log.Output(1, fmt.Sprintf("Creating a simulation from template %d for user %s", id, username))

	// Ask the server to create the clone and tell us the simulation id
	var result CloneResult
	body, err := auth.ProtectedResourceServerRequest(username, " create simulation ", `users/clone/`+t)
	if err != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Failed to complete clone because of %v", err))
		return
	}

	// read the simulation id
	jsonErr := json.Unmarshal(body, &result)
	if jsonErr != nil {
		utils.DisplayError(ctx, fmt.Sprintf("Couldn't decode the simulation id because of this error:%v", jsonErr))
		return
	} else {
		log.Output(1, fmt.Sprintf("Setting current simulation to be %d", result.Simulation_id))
		models.Users[username].CurrentSimulation = result.Simulation_id
	}

	models.Users[username].Initialize() // Wipe past history and create a new history record

	// Diagnostic - comment or uncomment as needed
	// s, _ := json.MarshalIndent(models.Users[username], "  ", "  ")
	// fmt.Printf("User record after creating the simulation is %s\n", string(s))

	if !api.FetchUserObjects(ctx, username) {
		utils.DisplayError(ctx, "Warning: we created this simulation but failed to retrieve all the data from the server")
	}
	ShowIndexPage(ctx)
}
