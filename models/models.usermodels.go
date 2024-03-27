// usergo
// User models and related objects

package models

import (
	"encoding/json"
)

const (
	Quantity = iota
	Value
	Price
)

// Full details of a user
// NOTE we do not store the password - this is handled by the remote server
// TODO factory function to initialize a new user
type UserData struct {
	Token             string              // The access token to use when requesting access to protected resources
	UserName          string              // Repeats the key in the map, which makes it easier to place in the admin dashboard template
	CurrentSimulation int                 // the id of the simulation that this user is currently using
	UserMessage       *UserMessage        // store for messages to be displayed to the user when appropriate
	LoggedIn          bool                // Is this user logged in?
	LastVisitedPage   string              // Remember what the user was looking at (used when an action is requested)
	ViewedTimeStamp   int                 // Indexes the History field. Selects what the user is viewing
	History           map[int]HistoryItem // the history of the current simulation
}

func (u UserData) Contents() string {
	b, _ := json.MarshalIndent(u.History, "", " ")
	return string(b)
}

// wipe the history clean and create a single new HistoryItem
func (u UserData) ReInitialize() {
	u.History = make(map[int]HistoryItem)
	u.History[0] = HistoryItem{Time_stamp: 0}
}

// Wrappers for the various lists
// TODO implement all this with an interface

// Wrapper for the SimulationList
func (u UserData) Simulations() *[]Simulation {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	// return &u.History[u.ViewedTimeStamp].SimulationList //TODO why can't I do
	return &h.SimulationList
}

// Wrapper for the CommodityList
func (u UserData) Commodities() *[]Commodity {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.CommodityList
}

// Wrapper for the IndustryList
func (u UserData) Industries() *[]Industry {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.IndustryList
}

// Wrapper for the ClassList
func (u UserData) Classes() *[]Class {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.ClassList
}

// Wrapper for the IndustryStockList
func (u UserData) IndustryStocks() *[]Industry_Stock {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.IndustryStockList
}

// Wrapper for the ClassStockList
func (u UserData) ClassStocks() *[]Class_Stock {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.ClassStockList
}

// Wrapper for the TraceList
func (u UserData) Traces() *[]Trace {
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.TraceList
}

// Format of responses from the server for post requests
// Specifically (so far), login or register.
type ServerMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// Abbreviated user details as supplied by server
// Used (for example) to obtain the current simulation after cloning
type UserServerData struct {
	UserName          string `json:"username"`
	Is_superuser      bool   `json:"is_superuser"`
	CurrentSimulation int    `json:"current_simulation"`
	Id                int    `json:"id"`
	Is_logged_in      bool   `json:"is_logged_in"`
}

// Transitory variable used to pick up information about the user from the server
var UserServerItem UserServerData

// Messages to the user are stored here and should be displayed by the relevant page handler
type UserMessage struct {
	StatusCode int
	Message    string
}

// contains the details of every user's simulations and their status, accessed by username
var Users = make(map[string]*UserData)

// List of basic user data, for use by the administrator
var AdminUserList []UserData
