// usergo
// User models and related objects

package models

import "errors"

// Messages to the user are stored here and should be displayed by the relevant page handler
type UserMessage struct {
	StatusCode int
	Message    string
}

// Full details of a user
// NOTE we do not store the password - this is handled by the remote server
type UserDatum struct {
	Token             string        // The access token to use when requesting access to protected resources
	UserName          string        `json:"username"`           // Repeats the key in the map, which makes it easier to place in the admin dashboard template
	CurrentSimulation int           `json:"current_simulation"` // the id of the simulation that this user is currently using
	UserMessage       *UserMessage  // store for messages to be displayed to the user when appropriate
	LoggedIn          bool          `json:"is_logged_in"` // Is this user logged in?
	LastVisitedPage   string        // Remember what the user was looking at (used when an action is requested)
	ViewedTimeStamp   int           // Indexes the History field. Selects what the user is viewing
	IsSuperuser       bool          `json:"is_superuser"`
	History           []HistoryItem // the history of the current simulation
}

// Constructor for a standard initial UserDatum.
// containing an empty slice of HistoryItems.
//
// THIS SHOULD NEVER BE NULL. It just makes everything too complicated.
//
// If user has no current simulation, History should be present, but empty.
//
// This is best achieved by using this constructor
func NewUserDatum(username string) UserDatum {
	new_datum := UserDatum{
		UserName:          username,
		Token:             "",
		CurrentSimulation: 0,
		LoggedIn:          false,
		LastVisitedPage:   "",
		ViewedTimeStamp:   0,
		IsSuperuser:       false,
		History:           make([]HistoryItem, 0),
	}
	return new_datum
}

// Wipe the history clean, create a single new HistoryItem and put it
// in the User's History,
//
// NOTE this History MUST be present, or we will get a null pointer error.
// Therefore we test for this and return an error if UserDatum.History is nil.
//
// Later, we may try to preserve the Histories of all active simulations locally
// but basically, it's just as easy to fetch them from the server.
func (u *UserDatum) Initialize() error {
	if u.History == nil {
		return errors.New("programme error: this user has no History Item")
	}
	u.History = append(u.History, NewHistoryItem())
	return nil
}

// Wrappers for the various lists
// These all return nil if the user has no simulations as yet.
// TODO a lot of boilerplate here but hard to remove
// TODO because of the nil test problem.

// Simulations is a special case, because the dashboard displays
// a list of the user's simulations. But if the user has none
// we have to make up a fake list with nothing in it, or the
// app crashes when preparing the Template.
// TODO there must be a better way.
func (u UserDatum) Simulations() *[]Simulation {
	if len(u.History) == 0 {
		var fakeList []Simulation = make([]Simulation, 0)
		return &fakeList
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.SimulationList
}

// Wrapper for the CommodityList
func (u UserDatum) Commodities() *[]Commodity {
	if len(u.History) == 0 {
		return nil
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.CommodityList
}

// Wrapper for the IndustryList
func (u UserDatum) Industries() *[]Industry {
	if len(u.History) == 0 {
		return nil
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.IndustryList
}

// Wrapper for the ClassList
func (u UserDatum) Classes() *[]Class {
	if len(u.History) == 0 {
		return nil
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.ClassList
}

// Wrapper for the IndustryStockList
func (u UserDatum) IndustryStocks() *[]Industry_Stock {
	if len(u.History) == 0 {
		return nil
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.IndustryStockList
}

// Wrapper for the ClassStockList
func (u UserDatum) ClassStocks() *[]Class_Stock {
	if len(u.History) == 0 {
		return nil
	}
	var v int = u.ViewedTimeStamp
	var h HistoryItem = u.History[v]
	return &h.ClassStockList
}

// Wrapper for the TraceList
func (u UserDatum) Traces() *[]Trace {
	if len(u.History) == 0 {
		return nil
	}
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

// contains the details of every user's simulations and their status, accessed by username
var Users = make(map[string]*UserDatum)

// List of basic user data, for use by the administrator
var AdminUserList []UserDatum
