// display.objects.go
// handlers to display the objects of the simulation on the user's browser

package display

import (
	"capfront/api"
	"capfront/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// display all commodities in the current simulation
func ShowCommodities(c *gin.Context) {
	c.HTML(http.StatusOK, "commodities.html", gin.H{
		"Title":       "Commodities",
		"commodities": models.CommodityList,
	})
}

// display all industries in the current simulation
func ShowIndustries(c *gin.Context) {
	c.HTML(http.StatusOK, "industries.html", gin.H{
		"Title":      "Industries",
		"industries": models.IndustryList,
	})
}

// display all classes in the current simulation
func ShowClasses(c *gin.Context) {
	c.HTML(http.StatusOK, "classes.html", gin.H{
		"Title":   "Classes",
		"classes": models.ClassList,
	})
}

// display all stocks in the current simulation
func ShowStocks(c *gin.Context) {
	c.HTML(http.StatusOK, "stocks.html", gin.H{
		"Title":  "Stocks",
		"stocks": models.StockList,
	})
}

// Display one specific commodity
func ShowCommodity(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.CommodityList); i++ {
		if id == models.CommodityList[i].Id {
			c.HTML(http.StatusOK, "commodity.html", gin.H{
				"Title":     "Commodity",
				"commodity": models.CommodityList[i],
			})
		}
	}
}

// Display one specific industry
func ShowIndustry(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.IndustryList); i++ {
		if id == models.IndustryList[i].Id {
			c.HTML(http.StatusOK, "industry.html", gin.H{
				"Title":    "Industry",
				"industry": models.IndustryList[i],
			})
		}
	}
}

// Display one specific class
func ShowClass(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.ClassList); i++ {
		if id == models.ClassList[i].Id {
			c.HTML(http.StatusOK, "class.html", gin.H{
				"Title": "Class",
				"class": models.ClassList[i],
			})
		}
	}
}

// Display one specific stock
func ShowStock(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id")) //TODO check user didn't do something stupid
	// TODO here and elsewhere create a method to get the simulation
	for i := 0; i < len(models.StockList); i++ {
		if id == models.StockList[i].Id {
			c.HTML(http.StatusOK, "stock.html", gin.H{
				"Title": "Stock",
				"stock": models.StockList[i],
			})
		}
	}
}

// services request to refresh all data from the server
// TODO this is probably not necessary; included for now for legacy and test purposes
func RefreshHandler(c *gin.Context) {
	api.Refresh()
	ShowIndexPage(c)
}

// Displays messages, and the economy
func ShowIndexPage(c *gin.Context) {
	// TODO implement display options
	var displayOptions = []string{
		`Values`,
		`Prices`,
		`Quantities`,
	}
	api.UserMessage = `This is the home page`
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":          "Economy",
		"industries":     models.IndustryList,
		"commodities":    models.CommodityList,
		"Message":        api.UserMessage,
		"DisplayOptions": displayOptions,
		"classes":        models.ClassList,
	})
}

// Fetch the trace from the local database
func ShowTrace(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"trace.html",
		gin.H{
			"Title": "Simulation Trace",
			"trace": models.TraceList,
		},
	)
}

// TODO the simulations should be confined to those belonging to the logged in user; the templates should be everything on offer
// retrieve all simulations from the local database and display them in the user dashboard
func UserDashboard(c *gin.Context) {

	api.FetchAPI(&api.ApiList[1]) //TODO this should be a map, indexed by the name of what we want. As a botch, we have put it in this location in
	// TODO deal with failure to return anything which is a programming error - this simulation should always be present, but
	// TODO (1) API provider could fuck up
	// TODO (2) connection might be broken
	c.HTML(http.StatusOK, "user-dashboard.html", gin.H{
		"Title":       "Dashboard",
		"simulations": models.MySimulations,
		"templates":   models.SimulationList,
	})
}
