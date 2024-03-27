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

	// Uncomment for more detailed diagnostics
	// _, file, no, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Printf("userStatus was called from %s#%d\n", file, no)
	// }

	// find out what the browser knows
	username, err := auth.Get_current_user(ctx)
	if err != nil {
		log.Printf("The client browser knows nothing about user %s", username)
		return "unknown", false, err
	}

	// find out what the server knows
	// Create a UserDatum and load it from the server
	synched_user := models.NewUserDatum(username)
	body, err := auth.ProtectedResourceServerRequest(username, "Synchronise with server", `users/`+username)
	if err != nil {
		log.Printf("Could not get user %s's details because:\n%v\n", username, err)
		return username, false, err
	}

	// Decode what the server knows
	err = json.Unmarshal(body, &synched_user)
	if err != nil {
		// couldn't decode it
		log.Printf("The server failed to inform us about user %s", username)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return username, false, err
	}

	// Ask the server whether it accepts that the user is logged in
	log.Printf("The server knows about user %s - ask if we are logged in", username)

	if !synched_user.LoggedIn {
		log.Printf("User %s is not logged in at the server", username)
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return username, false, err
	}

	log.Printf("The server knows about user %s and says it is logged in", username)

	// We agree with the server that this user can log in.
	// Now synch with the server in case something changed
	{
		if models.Users[username].CurrentSimulation != synched_user.CurrentSimulation {
			log.Printf("We are out of synch. Server thinks our simulation is %d and client says it is %d",
				synched_user.CurrentSimulation,
				models.Users[username].CurrentSimulation)
			if !api.FetchUserObjects(ctx, username) {
				log.Printf("We don't have a token. Redirecting to login")
				ctx.Redirect(http.StatusMovedPermanently, "/login")
				return username, false, nil
			}
		}

		models.Users[username].LastVisitedPage = ctx.Request.URL.Path
		models.Users[username].CurrentSimulation = synched_user.CurrentSimulation
		loginStatus = models.Users[username].LoggedIn
		return username, loginStatus, err
	}
}

// helper function to obtain the state of the current simulation
// if no user is logged in, return null state
func get_current_state(username string) string {
	this_user := models.Users[username]
	if this_user == nil {
		return "NO KNOWN USER"
	}
	this_user_history := this_user.History
	if len(this_user_history) == 0 {
		// User doesn't yet have any simulations
		return "NO SIMULATION YET"
	}

	id := this_user.CurrentSimulation
	sims := *this_user.Simulations()
	if sims == nil {
		return "UNKNOWN"
	}

	for i := 0; i < len(sims); i++ {
		s := sims[i]
		if s.Id == id {
			return s.State
		}
	}
	return "UNKNOWN"
}

// helper function to set the state of the current simulation
// if we fail it's a programme error so we don't test for that
func set_current_state(username string, new_state string) {
	this_user := models.Users[username]
	id := this_user.CurrentSimulation
	sims := *this_user.Simulations()
	log.Output(1, fmt.Sprintf("resetting state to %s for user %s", new_state, this_user.UserName))
	for i := 0; i < len(sims); i++ {
		s := &sims[i]
		if (*s).Id == id {
			(*s).State = new_state
			return
		}
		log.Output(1, fmt.Sprintf("simulation with id %d not found", id))
	}
}

// display all commodities in the current simulation
// use the cookie, which comes in the response, to identify the user
func ShowCommodities(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	state := get_current_state(username)

	ctx.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":          "Commodities",
		"commodities":    models.Users[username].Commodities(),
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all industries in the current simulation
func ShowIndustries(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":          "Industries",
		"industries":     models.Users[username].Industries(),
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all classes in the current simulation
func ShowClasses(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	state := get_current_state(username)
	ctx.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":          "Classes",
		"classes":        models.Users[username].Classes(),
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// Display one specific commodity
func ShowCommodity(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id"))
	// TODO here and elsewhere create a method to get the simulation
	// id := this_user.CurrentSimulation
	// sims := *this_user.Simulations()

	clist := *models.Users[username].Commodities()
	for i := 0; i < len(clist); i++ {
		if id == clist[i].Id {
			ctx.HTML(http.StatusOK, "commodity.html", gin.H{
				"Title":          "Commodity",
				"commodity":      clist[i],
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
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	ilist := *models.Users[username].Industries()
	for i := 0; i < len(ilist); i++ {
		if id == ilist[i].Id {
			ctx.HTML(http.StatusOK, "industry.html", gin.H{
				"Title":          "Industry",
				"industry":       ilist[i],
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
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	id, _ := strconv.Atoi(ctx.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	list := *models.Users[username].Classes()

	for i := 0; i < len(list); i++ {
		if id == list[i].Id {
			ctx.HTML(http.StatusOK, "class.html", gin.H{
				"Title":          "Class",
				"class":          list[i],
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
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	state := get_current_state(username)

	api.UserMessage = `This is the home page`

	clist := *models.Users[username].Commodities()
	ilist := *models.Users[username].Industries()
	cllist := *models.Users[username].Classes()

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Title":          "Economy",
		"industries":     ilist,
		"commodities":    clist,
		"Message":        models.Users[username].UserMessage.Message,
		"classes":        cllist,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// Fetch the trace from the local database
func ShowTrace(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	tlist := *models.Users[username].Traces()

	ctx.HTML(
		http.StatusOK,
		"trace.html",
		gin.H{
			"Title":          "Simulation Trace",
			"trace":          tlist,
			"username":       username,
			"loggedinstatus": loginStatus,
			"state":          state,
		},
	)
}

// Retrieve all templates, and all simulations belonging to this user, from the local database
// Display them in the user dashboard
func UserDashboard(ctx *gin.Context) {

	// Uncomment for more detailed diagnostics
	// _, file, no, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Printf("User Dashboard was called from %s#%d\n", file, no)
	// }

	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	slist := *models.Users[username].Simulations()

	ctx.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":          "Dashboard",
		"simulations":    slist,
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
		ctx.Redirect(http.StatusMovedPermanently, "/login")
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
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to delete simulation %d", username, id))
	auth.ProtectedResourceServerRequest(username, "Delete simulation", "simulations/delete/"+ctx.Param("id"))
	api.FetchUserObjects(ctx, username)
	UserDashboard(ctx)
}

func RestartSimulation(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Output(1, fmt.Sprintf("User %s wants to restart simulation %d", username, id))
	ctx.HTML(http.StatusOK, "notready.html", gin.H{
		"Title": "Not Ready",
	})
}

// display all industry stocks in the current simulation
func ShowIndustryStocks(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	islist := *models.Users[username].IndustryStocks()

	ctx.HTML(http.StatusOK, "industry_stocks.html", gin.H{
		"Title":          "Industry Stocks",
		"stocks":         islist,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}

// display all the class stocks in the current simulation
func ShowClassStocks(ctx *gin.Context) {
	username, loginStatus, _ := userStatus(ctx)
	if !loginStatus {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
		return
	}

	state := get_current_state(username)
	cslist := *models.Users[username].ClassStocks()

	ctx.HTML(http.StatusOK, "class_stocks.html", gin.H{
		"Title":          "Class Stocks",
		"stocks":         cslist,
		"username":       username,
		"loggedinstatus": loginStatus,
		"state":          state,
	})
}
