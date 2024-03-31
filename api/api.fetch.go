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

// Method of ApiItem which wraps the address of the user's local storage area.
//
//	item: provides the url that retrieves this table from the remote server.
//	username: which user should receive it.
//
// Returns nil if the List does not exist.
// This happens if the user is not currently in the client userList.
// or if the user has no simulations.
func (item ApiItem) Target(username string) any {
	u, ok1 := models.Users[username]

	if !ok1 {
		// This user is not in the local database
		fmt.Printf("Target reporting. User %v does not have a local user record", username)
		return nil
	}

	fmt.Printf("User %v is in the local database", u)
	// Diagnostics. Comment or uncomment as needed.
	fmt.Printf("Target reporting. User %v has a local user record\n", username)

	// Two special cases.
	// The TemplateList or the AdminUserList are both global lists
	// maintained by the admin user.
	if item.Name == `template` {
		return &models.TemplateList
	}
	if item.Name == `users` {
		return &models.AdminUserList
	}

	// Eventually we want to implement the code below. But we couldn't get it to work.
	// As a forensic step, we will put the switch code into FetchApi and see what
	// goes wrong.
	return nil
	// If this is a normal user, the data is stored in that user's first History record.

	// h := u.History[0]
	// fmt.Printf("Locating the target for a data fetch from the server, in history item %v\n", h)
	// switch item.Name {
	// case `simulation`:
	// 	return &h.SimulationList
	// case `commodity`:
	// 	return &h.CommodityList
	// case `industry`:
	// 	return &h.IndustryList
	// case `class`:
	// 	return &h.ClassList
	// case `industry_stock`:
	// 	return &h.IndustryStockList
	// case `class_stock`:
	// 	return &h.ClassStockList
	// case `trace`:
	// 	return &h.TraceList
	// default:
	// 	log.Output(1, fmt.Sprintf("Unknown dataset%s ", item.Name))
	// 	return nil
	// }
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
func FetchAdminObjects() bool {
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
	// Uncomment for more detailed diagnostics
	_, file, no, ok := runtime.Caller(1)
	if ok {
		fmt.Printf(" Fetch user objects was called from %s#%d\n", file, no)
	}
	for i := 0; i < len(ApiList); i++ {
		a := ApiList[i]
		fmt.Printf(" FetchUserObjects is fetching API item %d with name %s from URL %s\n", i, a.Name, a.ApiUrl)
		if !FetchAPI(&a, username) {
			fmt.Println("Cannot refresh from remote server; the user probably has no simulations. Do not continue ")
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
	// uncomment for more detailed diagnostics
	// _, file, no, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Printf("fetch API was called from %s#%d\n", file, no)
	// log.Output(1, fmt.Sprintf("User %s asked to fetch the table named %s from the URL %s\n", username, item.Name, item.ApiUrl))
	// }

	body, _ := auth.ProtectedResourceServerRequest(username, "Fetch Table", item.ApiUrl)
	// Check for an empty result. Not an error but tells us the user has no simulations yet.
	if body[0] == 91 && body[1] == 93 {
		log.Output(1, "INFORMATION: The server sent an empty table; this means the user has no simulations yet.")
	}

	// Uncomment for a little more diagnostic information
	// fmt.Printf("The server sent a table of length %d\n", len(string(body)))

	var jsonErr error
	// The next section is only used by the admin user.
	// It populates the global templates list.
	if item.Name == `template` {
		jsonErr = json.Unmarshal(body, &models.TemplateList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal the Template List because: %s", jsonErr))
			return false
		}
		return true
	}
	// The next section is only used by the admin user.
	// It populates the global user list.
	if item.Name == `users` {
		jsonErr = json.Unmarshal(body, &models.AdminUserList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal the User List because: %s", jsonErr))
			return false
		}
		return true
	}

	// It's a normal user, asking for its data
	// This section to be discarded in favour of using the Target method above.
	// Once we know how to finagle the address arithmetic

	u := models.Users[username]
	if len(u.History) == 0 {
		log.Output(1, fmt.Sprintf("User %s with no simulations asked for data", username))
		return false
	}
	// Uncomment for more detailed diagnostics
	// fmt.Printf("Unmarshalling data into history item %v\n", u.History[0])
	switch item.Name {
	case `simulation`:
		// jsonErr = json.Unmarshal(body, &h.SimulationList) // Does not work, find out why
		var sim []models.Simulation
		jsonErr = json.Unmarshal(body, &sim)
		models.Users[username].History[0].SimulationList = sim
	case `commodity`:
		// jsonErr = json.Unmarshal(body, &h.CommodityList) // Does not work, find out why ... etc
		var com []models.Commodity
		json.Unmarshal(body, &com)
		models.Users[username].History[0].CommodityList = com
	case `industry`:
		var ind []models.Industry
		json.Unmarshal(body, &ind)
		models.Users[username].History[0].IndustryList = ind
	case `class`:
		var cls []models.Class
		json.Unmarshal(body, &cls)
		models.Users[username].History[0].ClassList = cls
	case `industry_stock`:
		var is []models.Industry_Stock
		json.Unmarshal(body, &is)
		models.Users[username].History[0].IndustryStockList = is
	case `class_stock`:
		var cs []models.Class_Stock
		json.Unmarshal(body, &cs)
		models.Users[username].History[0].ClassStockList = cs
	case `trace`:
		var tra []models.Trace
		json.Unmarshal(body, &tra)
		models.Users[username].History[0].TraceList = tra
	default:
		log.Output(1, fmt.Sprintf("Unknown dataset%s ", item.Name))
		return false
	}
	if jsonErr != nil {
		log.Output(1, fmt.Sprintf("Failed to unmarshal template json because: %s", jsonErr))
		return false
	}

	// uncomment for verbose diagnostics
	// s, _ := json.MarshalIndent(models.Users[username], "  ", "  ")
	// fmt.Printf("User record after creating the simulation is %s\n", string(s))
	return true
}
