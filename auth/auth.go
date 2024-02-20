// auth.auth.go
// handles login, registration, authentication

package auth

import (
	"fmt"
	"net/http"
	"strings"
)

type Token struct {
	Token string `json:"token"`
}

var UserToken Token
var LoggedInUser string

//TODO 18.02.2024 as far as I can see this doesn't actually do anything useful, even if it ever did, except for commenting critically on the general goings-on

func VerifyTokenController(w http.ResponseWriter, r *http.Request) {
	prefix := "Bearer "
	authHeader := r.Header.Get("Authorization")
	reqToken := strings.TrimPrefix(authHeader, prefix)

	if authHeader == "" || reqToken == authHeader {
		fmt.Println("Authentication header not present or malformed")
		return
	}

	// fmt.Println("Token is valid")
}
