// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package display

import (
	"capfront/api"
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// helper function for most display handlers
// retrieves a user cookie using Get_current_user and extracts the login status
// for convenience, returns the username, whether the user is logged in, and any error
// sets 'LastVisitedPage' so we can return here after an action
func userStatus(ctx *gin.Context) (string, bool, error) {
	var loginStatus bool = false

	// find out what the browser knows

	username, err := auth.Get_current_user(ctx)
	if err != nil {
		log.Printf("The client browser knows nothing about user %s", username)
		return "unknown", false, err
	} else {

		// find out what the server knows

		synched_user := models.UserServerData{}
		body, err := auth.ProtectedResourceServerRequest(username, "Synchronise with server", `users/`+username)
		if err != nil {
			log.Printf("The server knows nothing about user %s", username)
			return username, false, err
		}

		err = json.Unmarshal(body, &synched_user)

		// the server knows something

		if err != nil {
			log.Printf("The server failed to inform us about user %s", username)
			ctx.Redirect(http.StatusFound, "/login")
			// TODO tell the user why she is being asked to log in again
			return username, false, err
		} else

		// Ask the server whether it accepts that the user is logged in

		if !synched_user.Is_logged_in {
			log.Printf("User %s is not logged in at the server", username)
			ctx.Redirect(http.StatusFound, "/login")
			// TODO tell the user why she is being asked to log in again
			return username, false, err
		}

		// We agree with the server that this user can log in.
		// Now synch with the server in case something changed
		{
			if models.Users[username].CurrentSimulation != synched_user.CurrentSimulation {
				log.Printf("We are out of synch. Server thinks our simulation is %d and client says it is %d",
					synched_user.CurrentSimulation,
					models.Users[username].CurrentSimulation)
				api.Refresh(ctx, username)
			}

			models.Users[username].LastVisitedPage = ctx.Request.URL.Path
			models.Users[username].CurrentSimulation = synched_user.CurrentSimulation
			loginStatus = models.Users[username].LoggedIn
			return username, loginStatus, err
		}
	}
}

// helper function to obtain the state of the current simulation
// if no user is logged in, return null state
func get_current_state(username string) string {
	this_user := models.Users[username]
	if this_user == nil {
		return "NO SIMULATION YET"
	}
	this_simulation_id := this_user.CurrentSimulation
	for i := 0; i < len(this_user.SimulationList); i++ {
		s := this_user.SimulationList[i]
		if s.Id == this_simulation_id {
			return s.State
		}
	}
	return "UNKNOWN"
}

// helper function to set the state of the current simulation
// if we fail it's a programme error so we don't test for that
func set_current_state(username string, new_state string) {
	this_user := models.Users[username]
	this_simulation_id := this_user.CurrentSimulation
	log.Output(1, fmt.Sprintf("resetting state to %s for user %s", new_state, this_user.UserName))
	for i := 0; i < len(this_user.SimulationList); i++ {
		s := &this_user.SimulationList[i]
		if (*s).Id == this_simulation_id {
			(*s).State = new_state
			return
		}
		log.Output(1, fmt.Sprintf("simulation with id %d not found", this_simulation_id))
	}
}

// display all commodities in the current simulation
// use the cookie, which comes in the response, to identify the user
func ShowCommodities(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	state := get_current_state(username)

	ctx.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":          "Commodities",
		"commodities":    models.Users[username].CommodityList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all industries in the current simulation
func ShowIndustries(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":          "Industries",
		"industries":     models.Users[username].IndustryList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all classes in the current simulation
func ShowClasses(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":          "Classes",
		"classes":        models.Users[username].ClassList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// Display one specific commodity
func ShowCommodity(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id"))
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.Users[username].CommodityList); i++ {
		if id == models.Users[username].CommodityList[i].Id {
			ctx.HTML(http.StatusOK, "commodity.html", gin.H{
				"Title":          "Commodity",
				"commodity":      models.Users[username].CommodityList[i],
				"username":       username,
				"loggedinstatus": loginStatus,
				"state":          state,
			})
		}
	}
}

// Display one specific industry
func ShowIndustry(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.Users[username].IndustryList); i++ {
		if id == models.Users[username].IndustryList[i].Id {
			ctx.HTML(http.StatusOK, "industry.html", gin.H{
				"Title":          "Industry",
				"industry":       models.Users[username].IndustryList[i],
				"username":       username,
				"loggedinstatus": loginStatus,
				"state":          state,
			})
		}
	}
}

// Display one specific class
func ShowClass(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.Users[username].ClassList); i++ {
		if id == models.Users[username].ClassList[i].Id {
			ctx.HTML(http.StatusOK, "class.html", gin.H{
				"Title":          "Class",
				"class":          models.Users[username].ClassList[i],
				"username":       username,
				"loggedinstatus": loginStatus,
				"state":          state,
			})
		}
	}
}

// Displays snapshot of the economy
// TODO parameterise the templates to reduce boilerplate
func ShowIndexPage(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}
	state := get_current_state(username)

	api.UserMessage = `This is the home page`
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title":          "Economy",
		"industries":     models.Users[username].IndustryList,
		"commodities":    models.Users[username].CommodityList,
		"Message":        models.Users[username].UserMessage.Message,
		"DisplayOptions": models.Quantity,
		"classes":        models.Users[username].ClassList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// Fetch the trace from the local database
func ShowTrace(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(
		http.StatusOK,
		"trace.html",
		gin.H{
			"Title":          "Simulation Trace",
			"trace":          models.Users[username].TraceList,
			"username":       username,
			"loggedinstatus": loginStatus,
			"state":          state,
		},
	)
}

// Retrieve all templates, and all simulations belonging to this user, from the local database
// Display them in the user dashboard
func UserDashboard(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":          "Dashboard",
		"simulations":    models.Users[username].SimulationList,
		"templates":      models.TemplateList,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// a diagnostic endpoint to display the data in the system
func DataHandler(ctx *gin.Context) {
	// username, loginStatus, _ := userStatus(ctx)
	// b, err := json.Marshal(models.Users)
	// if err != nil {
	// 	fmt.Println("Could not marshal the Users object")
	// 	return
	// }
	ctx.JSON(http.StatusOK, models.Users)
}

func SwitchSimulation(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to switch to simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

func DeleteSimulation(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to delete simulation %d", username, id))
	auth.ProtectedResourceServerRequest(username, "Delete simulation", "simulations/delete/"+ctx.Param("id"))
	api.Refresh(ctx, username)
	UserDashboard(ctx)
}

func RestartSimulation(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusFound, "/login")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to restart simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}
