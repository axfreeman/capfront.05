package main

import (
	"capfront/api"
	"capfront/auth"
	"capfront/display"
	"capfront/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// Runs once at startup
// Loads environment variables
// Downloads user details from server
// Downloads starter templates
func Initialise() {
	// err := gotdotenv.Load()                 // 👈 load .env file
	auth.SECRET_ADMIN_PASSWORD = "insecure" // TODO get this from settings file

	// At startup, create an admin user record and authenticate it with the server.
	// We assume there is an admin user on the server with the password used here.
	admin_user := models.NewUserDatum(`admin`)
	models.Users["admin"] = &admin_user
	token, result := display.ServerLogin("admin", auth.SECRET_ADMIN_PASSWORD)

	if token == nil {
		log.Fatalf("Server failed at startup. It said:\n%v", result)
	}

	admin_user.Token = token.(string)

	// Get a list of templates from the server and put it in TemplateList
	// Get a list of users from the server and put it in AdminUserList
	if !api.FetchAdminObjects() {
		log.Fatal("Could not retrieve enough information from the server. Stopping")
	}
}

func main() {
	r := gin.New()
	r.LoadHTMLGlob("./templates/**/*") // load all the templates in the templates folder
	fmt.Println("Welcome to capitalism")
	r.GET("/action/:action", display.ActionHandler)
	r.GET("/commodities", display.ShowCommodities)
	r.GET("/industries", display.ShowIndustries)
	r.GET("/classes", display.ShowClasses)
	r.GET("/industry_stocks", display.ShowIndustryStocks)
	r.GET("/class_stocks", display.ShowClassStocks)
	r.GET("/industry/:id", display.ShowIndustry)
	r.GET("/commodity/:id", display.ShowCommodity)
	r.GET("/class/:id", display.ShowClass)
	r.GET("/trace", display.ShowTrace)
	r.GET("/admin/dashboard", display.AdminDashboard)
	r.GET("/admin/reset", display.AdminReset)
	r.GET("/login", display.CaptureLoginRequest)
	r.POST("/user/login", display.HandleLoginRequest)
	r.GET("/logout", display.ClientLogoutRequest)
	r.GET("/register", display.CaptureRegisterRequest)
	r.POST("/user/register", display.HandleRegisterRequest)
	r.GET("/user/create/:id", display.CreateSimulation)
	r.GET("/user/dashboard", display.UserDashboard)
	r.GET("/user/switch/:id", display.SwitchSimulation)
	r.GET("/user/delete/:id", display.DeleteSimulation)
	r.GET("/user/restart/:id", display.RestartSimulation)
	r.GET("/index/", display.ShowIndexPage)
	r.GET("/data/", display.DataHandler)
	r.GET("/displaymode", display.DisplayMode)
	r.GET("/", display.ShowIndexPage)
	Initialise()
	r.Run() // Run the server

}
