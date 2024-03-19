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
	admin_user := models.UserData{LoggedIn: false, UserName: "admin", Token: ""}
	models.Users["admin"] = &admin_user
	serverPayload, err := display.ServerLogin("admin", auth.SECRET_ADMIN_PASSWORD)

	if err != nil {
		log.Fatalf("Server failed at startup. It said:\n%v", serverPayload["message"])
	}

	api.FetchAPI(&api.ApiList[0], "admin") // get templates
	api.FetchAPI(&api.ApiList[1], "admin") // get user details
	// Copy the list we just downloaded into the UserList
	// Can probably download directly into UserList
	// but I wasn't sure how the unMarshalling would affect the nested arrays
	for _, item := range models.AdminUserList {
		user := models.UserData{LoggedIn: false, UserName: item.UserName, Token: ""}
		models.Users[item.UserName] = &user
	}
	ListData()
}

// short diagnostic function to display user and template data
func ListData() {
	fmt.Printf("\nTemplateList has %d elements which are:\n", len(models.TemplateList))
	for i := 0; i < len(models.TemplateList); i++ {
		fmt.Println(models.TemplateList[i])
	}
	fmt.Printf("AdminUserList has %d elements which are:\n", len(models.AdminUserList))
	for i := 0; i < len(models.AdminUserList); i++ {
		fmt.Println(models.AdminUserList[i])
	}
}

func main() {
	r := gin.Default()
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
	r.GET("/logout", display.ClientLogoutRequest)
	r.POST("/user/login", display.ClientLoginRequest)
	r.GET("/user/create/:id", display.CreateSimulation)
	r.GET("/user/dashboard", display.UserDashboard)
	r.GET("/user/switch/:id", display.SwitchSimulation)
	r.GET("/user/delete/:id", display.DeleteSimulation)
	r.GET("/user/restart/:id", display.RestartSimulation)
	r.GET("/index/", display.ShowIndexPage)
	r.GET("/data/", display.DataHandler)
	r.GET("/register", display.Register)
	r.GET("/displaymode", display.DisplayMode)
	r.GET("/", display.ShowIndexPage)
	Initialise()
	r.Run() // Run the server

}
