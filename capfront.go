package main

import (
	"capfront/api"
	"capfront/display"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./templates/**/*") // load all the templates in the templates folder
	fmt.Println("Welcome to capitalism")
	r.GET("/commodities", display.ShowCommodities)
	r.GET("/industries", display.ShowIndustries)
	r.GET("/classes", display.ShowClasses)
	r.GET("/stocks", display.ShowStocks)
	r.GET("/industry/:id", display.ShowIndustry)
	r.GET("/commodity/:id", display.ShowCommodity)
	r.GET("/stock/:id", display.ShowStock)
	r.GET("/class/:id", display.ShowClass)
	r.GET("/trace", display.ShowTrace)
	r.GET("/action/:action", display.ActionHandler)
	r.GET("/login", display.CaptureLoginRequest)
	r.POST("/requestlogin", display.ServiceLoginRequest)
	r.GET("/user/dashboard", display.UserDashboard)
	r.GET("/refresh", display.RefreshHandler) // Fetch the data from the remote server
	r.GET("/", display.ShowIndexPage)
	r.GET("/user/select/:id", display.CreateSimulation)
	// TODO probably redundant
	// r.GET("/simulation/:id", GetSimulation)
	//TODO implement this
	// r.GET("/display/:mode", DisplayHandler)

	display.Fake_login() // temporary testing fix
	api.Refresh()
	r.Run() // Run the server

}
