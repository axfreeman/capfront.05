// api.fetch.go
// handlers to fetch objects from the remote server

package api

import (
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/gin-gonic/gin"
)

var UserMessage string

// Contains the information needed to fetch data for one model from the remote server.
// Name is a description, just for diagnostic purposes.
// ApiURL is the endpoint to get the data from the server.
type ApiItem struct {
	Name   string // the data to be obtained
	ApiUrl string // the url to be used in accessing the backend
}

// yields the Address of the List which stores the data that is fetched into ApiItem
// by a call to the server
func (a ApiItem) Target(username string) any {
	var hlist map[int]models.HistoryItem = models.Users[username].History
	var h models.HistoryItem = hlist[0]
	switch a.Name {
	case `template`:
		return &models.TemplateList
	case `users`:
		return &models.AdminUserList
	case `simulation`:
		return &h.SimulationList
	case `commodity`:
		return &h.CommodityList
	case `industry`:
		return &h.IndustryList
	case `class`:
		return &h.ClassList
	case `industry_stock`:
		return &h.IndustryStockList
	case `class_stock`:
		return &h.ClassStockList
	case `trace`:
		return &h.TraceList
	default:
		log.Output(1, fmt.Sprintf("Unknown dataset%s ", a.Name))
		return nil
	}
}

// a list of items needed to fetch data from the remote server
var ApiList = [7]ApiItem{
	{`simulation`, `simulations/mine`},
	{`commodity`, `commodities/`},
	{`industry`, `industries/`},
	{`class`, `classes/`},
	{`industry_stock`, `stocks/industry`},
	{`class_stock`, `stocks/class`},
	{`trace`, `trace/`},
}

// Populates the two global objects required by the simulation.
// These are the list of users, and the list of templates.
// Should only be called by the admin user.
// OR when a new user registers.
// Is also called at startup.
func FetchAdminObjects(ctx *gin.Context, username string) bool {
	if !FetchAPI(&ApiItem{`template`, `simulations/templates`}, `admin`) {
		return false
	}
	if !FetchAPI(&ApiItem{`users`, `users/`}, `admin`) {
		return false
	}
	return true
}

// Iterates through ApiList to refresh all user objects.
// Returns false if any table fails.
// Returns true if all tables succeed.
func FetchUserObjects(ctx *gin.Context, username string) bool {
	for i := 2; i < len(ApiList); i++ {
		a := ApiList[i]
		if !FetchAPI(&a, username) {
			fmt.Println("Cannot refresh from remote server; giving up")
			return false
		}
	}
	log.Output(1, "Refresh complete")
	return true
}

// fetch the data specified by item for user.
// if we got something, return true.
// if not, for whatever reason, return false.
func FetchAPI(item *ApiItem, username string) (result bool) {
	fmt.Printf("User %s asked to fetch the table named %s\n", username, item.ApiUrl)
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("fetch API was called from %s#%d\n", file, no)
	}

	body, _ := auth.ProtectedResourceServerRequest(username, "Fetch Table", item.ApiUrl)
	// Check for an empty result. Not necessarily an error but useful to know.
	if body[0] == 91 && body[1] == 93 {
		log.Output(1, "INFORMATION: The server sent an empty table")
	}

	// more detailed diagnostics. Comment in production version
	fmt.Printf("The server sent %s\n", string(body))

	var jsonErr error = json.Unmarshal(body, item.Target(username))

	if jsonErr != nil {
		log.Output(1, fmt.Sprintf("Failed to unmarshal template json because: %s", jsonErr))
		return false
	}

	// uncomment for verbose diagnostics
	fmt.Printf("After loading, the models map for user %s is:", username)
	fmt.Print(models.Users[username].History)
	fmt.Println("")
	return true
}

// diagnostic helper function
func PrintUsers() {
	b, err := json.MarshalIndent(models.Users, " ", " ")
	if err != nil {
		fmt.Println("Could not marshal the Users object")
		return
	}
	fmt.Println(string(b))
}
