// api.fetch.go
// handlers to fetch objects from the remote server

package api

import (
	"capfront/auth"
	"capfront/models"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// defines the API of the remote server
// TODO automate which API is used depending on whether running locally or in production
// var APIsource = `http://127.0.0.1:8000/`

var APIsource = `https://simcap-a09ebdf4b7cb.herokuapp.com/`

// TODO display this in a message section in the browser
// NOTE why can't this be in the main package? Should it be in some separate global package?
// Causes import cycle problems if it is placed in display package
// contains information to be displayed for the user to see, when appropriate
var UserMessage string

// Contains the information needed to fetch data for one model from the remote server
type ApiItem struct {
	Name     string      // the data to be obtained
	ApiUrl   string      // the url to be used in accessing the backend
	DataList interface{} // temporary store (tells json.Unmarshal what the fields are called and where to find them)
}

// a list of items needed to fetch data from the remote server
var ApiList = [7]ApiItem{
	{`template`, APIsource + `simulations/templates/`, &models.SimulationList},
	{`simulation`, APIsource + `simulations/mine/`, &models.SimulationList},
	{`commodity`, APIsource + `commodities/`, &models.CommodityList},
	{`industry`, APIsource + `industries/`, &models.IndustryList},
	{`class`, APIsource + `classes/`, &models.ClassList},
	{`stock`, APIsource + `stocks/`, &models.StockList},
	{`trace`, APIsource + `trace/`, &models.TraceList},
}

// refresh all objects from the remote server
// should be called whenever the server has undertaken an action which changes the contents of any object
func Refresh() {
	UserMessage = "Data refresh complete" // default response which will be overwritten if there is a failure
	for i := range ApiList {
		a := ApiList[i]
		// log.Output(1, fmt.Sprintf("Refreshing the %s table from server", a.Name))
		if !FetchAPI(&a) {
			UserMessage = "Cannot refresh from remote server; displaying existing stored data"
		}
		// TODO set the time stamp of the object where appropriate
	}
	log.Output(1, UserMessage) // TODO display this somewhere in the browser
}

// fetch the data specified by a
func FetchAPI(a *ApiItem) (result bool) {
	// log.Output(1, fmt.Sprintf("API call to get %s for display array %s", a.ApiUrl, a.Name))
	client := http.Client{Timeout: time.Second * 2} // Timeout after 2 seconds

	req, err := http.NewRequest(http.MethodGet, a.ApiUrl, nil)
	if err != nil {
		log.Output(1, fmt.Sprintf("Failed to generate a request for URL %s because of %v", a.ApiUrl, err))
		return false
	}

	req.Header.Set("User-Agent", "Capitalism reader")
	req.Header.Add("Authorization", "Bearer "+auth.UserToken.Token)
	auth.VerifyTokenController(nil, req)

	res, getErr := client.Do(req)
	if getErr != nil {
		log.Output(1, fmt.Sprintf("Failed to service a request for the URL %s because of %v", a.ApiUrl, getErr))
		return false
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Output(1, fmt.Sprintf("Failed to read the body of the request because: %s", readErr))
		return false
	}

	log.Output(1, fmt.Sprintf("Refreshed the %s table from API %s", a.Name, a.ApiUrl))
	// log.Output(1, fmt.Sprintf("The result was %s", body))
	// TODO deal with the exceptions below in a more organised way
	if string(body) == `{"detail":"Not authenticated"}` {
		log.Output(1, "The user is not authenticated")
		return false // User not logged in
	}
	if string(body) == `{"detail":"Signature has expired"}` {
		log.Output(1, "The user's signature has expired")
		return false // User logged in but signature has expired
	}

	switch a.Name {
	// TODO replace by methods and remove boilerplate
	case `template`:
		jsonErr := json.Unmarshal(body, &models.SimulationList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal template json because: %s", jsonErr))
		}
	case `simulation`:
		jsonErr := json.Unmarshal(body, &models.MySimulations)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal simulation json because: %s", jsonErr))
		}
	case `commodity`:
		jsonErr := json.Unmarshal(body, &models.CommodityList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal commodity json because: %s", jsonErr))
		}
	case `industry`:
		jsonErr := json.Unmarshal(body, &models.IndustryList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal commodity json because: %s", jsonErr))
		}
	case `class`:
		jsonErr := json.Unmarshal(body, &models.ClassList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal commodity json because: %s", jsonErr))
		}
	case `stock`:
		jsonErr := json.Unmarshal(body, &models.StockList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal commodity json because: %s", jsonErr))
		}
	case `trace`:
		jsonErr := json.Unmarshal(body, &models.TraceList)
		if jsonErr != nil {
			log.Output(1, fmt.Sprintf("Failed to unmarshal commodity json because: %s", jsonErr))
		}
	default:
		log.Output(1, fmt.Sprintf("Unknown dataset%s ", a.Name))
	}
	return true
}
