// api.fetch.go
// handlers to fetch objects from the remote server

package api

import (
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"log"

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
var ApiList = [10]ApiItem{
	{`template`, `simulations/templates`},
	{`users`, `users/`},
	{`simulation`, `simulations/mine`},
	{`commodity`, `commodities/`},
	{`industry`, `industries/`},
	{`class`, `classes/`},
	{`stock`, `stocks/`},
	{`industry_stock`, `stocks/industry`},
	{`class_stock`, `stocks/class`},
	{`trace`, `trace/`},
}

// Iterates through ApiList to refresh all objects owned by the user
// from the remote server, by invoking FetchAPI.
// TODO ctx is not needed.
func Refresh(ctx *gin.Context, username string) {
	for i := range ApiList {
		a := ApiList[i]
		if !FetchAPI(&a, username) {
			// If one fetch fails, there is no point continuing because
			// there has been a login failure or the server is down.
			// TODO handle this so the caller knows something went wrong.
			fmt.Println("Cannot refresh from remote server; giving up")
			return
		}
	}
	log.Output(1, "Refresh complete")
}

// fetch the data specified by item for user.
// if we got something, return true.
// if not, for whatever reason, return false.
func FetchAPI(item *ApiItem, username string) (result bool) {
	body, _ := auth.ProtectedResourceServerRequest(username, "Fetch Table", item.ApiUrl)
	// Check for an empty result.
	// This can happen, but we need to know it did.
	if string(body) == `[]` {
		fmt.Println("The result was an empty table")
		return false
	}
	var jsonErr error
	switch item.Name {
	case `template`:
		jsonErr = json.Unmarshal(body, &models.TemplateList)
	case `users`:
		jsonErr = json.Unmarshal(body, &models.AdminUserList)
	case `simulation`:
		jsonErr = json.Unmarshal(body, &models.Users[username].SimulationList)
	case `commodity`:
		jsonErr = json.Unmarshal(body, &models.Users[username].CommodityList)
	case `industry`:
		jsonErr = json.Unmarshal(body, &models.Users[username].IndustryList)
	case `class`:
		jsonErr = json.Unmarshal(body, &models.Users[username].ClassList)
	case `stock`:
		jsonErr = json.Unmarshal(body, &models.Users[username].StockList)
	case `industry_stock`:
		jsonErr = json.Unmarshal(body, &models.Users[username].IndustryStockList)
	case `class_stock`:
		jsonErr = json.Unmarshal(body, &models.Users[username].ClassStockList)
	case `trace`:
		jsonErr = json.Unmarshal(body, &models.Users[username].TraceList)
	default:
		log.Output(1, fmt.Sprintf("Unknown dataset%s ", item.Name))
	}
	if jsonErr != nil {
		log.Output(1, fmt.Sprintf("Failed to unmarshal template json because: %s", jsonErr))
		return false
	}

	// uncomment for verbose diagnostics
	// fmt.Println("After loading, the models map for user guest is:")
	// PrintUsers()
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
