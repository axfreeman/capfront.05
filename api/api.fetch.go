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

// Contains the information needed to fetch data for one model from the remote server
// TODO use interfacing to add a destination field
type ApiItem struct {
	Name   string // the data to be obtained
	ApiUrl string // the url to be used in accessing the backend
}

// defines the API of the remote server
var UserMessage string

// a list of items needed to fetch data from the remote server
var ApiList = [9]ApiItem{
	{`template`, `simulations/templates`},
	{`users`, `users/`},
	{`simulation`, `simulations/mine`},
	{`commodity`, `commodities/`},
	{`industry`, `industries/`},
	{`class`, `classes/`},
	{`industry_stock`, `stocks/industry`},
	{`class_stock`, `stocks/class`},
	{`trace`, `trace/`},
}

// Iterates through ApiList to refresh all user objects.
// Omits the first two items in the table which are reserved for admin user
// Returns false if any table fails.
// Returns true if all tables succeed.
func Refresh(ctx *gin.Context, username string) bool {
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
	var jsonErr error

	var v int = models.Users[username].ViewedTimeStamp
	var hlist = models.Users[username].History

	// TODO this is a bodge. We need some separate method to get the Template and Userlists
	if username == `admin` {
		// The admin wants to populate the global template and user lists
		switch item.Name {
		case `template`:
			jsonErr = json.Unmarshal(body, &models.TemplateList)
		case `users`:
			jsonErr = json.Unmarshal(body, &models.AdminUserList)
		case `default`:
			log.Fatal("Admin user did something it shouldn't")
		}
	} else {
		h, ok := hlist[v]
		// Check for programme error: if h does not exist, we did something wrong.
		if !ok {
			log.Output(1, "The application tried to populate a non-existent history item")
			return
		}
		fmt.Printf("There are %d items in the History Map and the current Time Stamp is %d\n", len(hlist), v)
		switch item.Name {
		case `template`:
			// THIS SHOULD NOT HAPPEN! But if it does, report it
			log.Output(1, fmt.Sprintf("Non-admin user %s tried to populat the template list", username))
			jsonErr = json.Unmarshal(body, &models.TemplateList)
		case `users`:
			// THIS SHOULD NOT HAPPEN! But if it does, report it
			log.Output(1, fmt.Sprintf("Non-admin user %s tried to populat the template list", username))
			jsonErr = json.Unmarshal(body, &models.AdminUserList)
		case `simulation`:
			jsonErr = json.Unmarshal(body, &h.SimulationList)
		case `commodity`:
			jsonErr = json.Unmarshal(body, &h.CommodityList)
		case `industry`:
			jsonErr = json.Unmarshal(body, &h.IndustryList)
		case `class`:
			jsonErr = json.Unmarshal(body, &h.ClassList)
		case `industry_stock`:
			jsonErr = json.Unmarshal(body, &h.IndustryStockList)
		case `class_stock`:
			jsonErr = json.Unmarshal(body, &h.ClassStockList)
		case `trace`:
			jsonErr = json.Unmarshal(body, &h.TraceList)
		default:
			log.Output(1, fmt.Sprintf("Unknown dataset%s ", item.Name))
			return false
		}

	}
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
