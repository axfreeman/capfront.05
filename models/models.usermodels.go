// usergo
// User models and related objects

package models

const (
	Quantity = iota
	Value
	Price
)

// Full details of a user
// NOTE we do not store the password - this is handled by the remote server
type UserData struct {
	Token             string       // The access token to use when requesting access to protected resources
	UserName          string       // Repeats the key in the map, which makes it easier to place in the admin dashboard template
	CurrentSimulation int          // the id of the simulation that this user is currently using
	UserMessage       *UserMessage // store for messages to be displayed to the user when appropriate
	LoggedIn          bool         // Is this user logged in?
	LastVisitedPage   string       // Remember what the user was looking at (used when an action is requested)
	DisplayOption     string       // price, value or size TODO make this type-safe? Probably overkill
	SimulationList    []Simulation // all the simulations this user has created
	CommodityList     []Commodity  // all the commodity objects this user has created
	IndustryList      []Industry   // ...
	ClassList         []Class
	StockList         []Stock
	TraceList         []Trace
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
