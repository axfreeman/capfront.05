// display.actions.go
// handlers for actions requested by the user

package display

import (
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// dispatches the action requested in the URL (eg /actions/trade will perform the trade action, etc)
type Action struct {
	A string `uri:"action"`
}

// Handles requests for the server to take an action comprising a stage
// of the circuit (demand,supply, trade, produce, invest), corresponding
// to a button press. This is specified by the URL parameter 'act'
// Having requested the action from ths server, sets 'state' to the next
// stage of the circuit and redisplays whatever the user was looking at
func ActionHandler(ctx *gin.Context) {
	log.Output(1, "Entered actionHandler")
	var param Action
	err := ctx.ShouldBindUri(&param)
	if err != nil {
		fmt.Println("Something went wrong", err)
		ctx.String(http.StatusBadRequest, "Malformed URL")
		return
	}
	act := ctx.Param("action")
	username, _ := auth.Get_current_user(ctx)
	lastVisitedPage := models.Users[username].LastVisitedPage
	log.Output(1, fmt.Sprintf("User %s wants the server to do %s\n", username, act))
	log.Output(1, fmt.Sprintf("Last visited page %s", lastVisitedPage))
	auth.ProtectedResourceServerRequest(username, act, `action/`+act)
	api.Refresh(ctx, username)
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
	v := visitedPageURL[0]
	if v == `commodities` || v == `industries` || v == `classes` || v == `stocks` {
		ctx.Redirect(http.StatusFound, lastVisitedPage)
	} else {
		ctx.Redirect(http.StatusFound, "/index")
	}
	// //TODO set time stamp
	// TODO why tf don't the log statements (or their println equivalents) produce any output?
}

// Creates a new simulation for the logged-in user, from the template specified by the 'id' parameter
func CreateSimulation(ctx *gin.Context) {
	username, _ := auth.Get_current_user(ctx)
	template_id := ctx.Param("id")
	body, _ := auth.ProtectedResourceServerRequest(username, " get user details ", `users/`+username)
	jsonErr := json.Unmarshal(body, &models.UserServerItem)
	auth.ProtectedResourceServerRequest(username, " create simulation ", `users/clone/`+template_id)
	if jsonErr != nil {
		log.Output(1, "Failed to obtain user details while creating a new simulation - cannot set current simulation right now")
	} else {
		log.Output(1, fmt.Sprintf("Setting current simulation to be %d", models.UserServerItem.CurrentSimulation))
		models.Users[username].CurrentSimulation = models.UserServerItem.CurrentSimulation
	}
	api.Refresh(ctx, username)
	ShowIndexPage(ctx)
}
