// methods.simulation.go
// class methods of the objects specified in models.simulation.go
// TODO use address arithmetic to access all lists
package models

import (
	"strconv"
)

// TODO eliminate boilerplate by making generic
// TODO see https://github.com/jose78/go-collection/blob/master/collections/collection.go  for suggestions

//METHODS OF INDUSTRIES

// crude searches without database implementation
// justified because we can avoid the complications of a database implementation
// and the size of the tables is not large, because they are provided on a per-user basis
// However as the simulations get large, this may become more problematic (let's find out pragmatically)
// In that case some more sophisticated system, such as a local database, may be needed
// A simple solution would be to add direct links to related objects in the models
// perhaps populated by an asynchronous process in the background

// A default stock returned if any condition is not met (that is, if the predicated stock does not exist)
// Used to signal to the user that there has been a programme error
var NotFoundStock = Stock{
	Id:            0,
	Simulation_id: 0,
	Time_Stamp:    0,
	Owner_id:      0,
	Commodity_id:  0,
	Name:          "NOT FOUND",
	Owner_type:    "PROGRAMME ERROR",
	Usage_type:    "PROGRAMME ERROR",
	Size:          -1,
	Value:         -1,
	Price:         -1,
	Requirement:   -1,
	Demand:        -1,
}

var NotFoundCommodity = Commodity{
	Id:                          0,
	Name:                        "NOT FOUND",
	Simulation_id:               0,
	Time_Stamp:                  0,
	Origin:                      "UNDEFINED",
	Usage:                       "UNDEFINED",
	Size:                        0,
	Total_Value:                 0,
	Total_Price:                 0,
	Unit_Value:                  0,
	Unit_Price:                  0,
	Turnover_Time:               0,
	Demand:                      0,
	Supply:                      0,
	Allocation_Ratio:            0,
	Display_Order:               0,
	Image_Name:                  "UNDEFINED",
	Tooltip:                     "UNDEFINED",
	Monetarily_Effective_Demand: 0,
	Investment_Proportion:       0,
}

// returns the money stock of the given industry
// WAS err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", industry.Id, "Money").StructScan(&stock)
func (industry Industry) MoneyStock() Stock {
	username := industry.UserName
	stockList := (Users[username].StockList)
	for i := 0; i < len(stockList); i++ {
		s := stockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Money`) {
			return s
		}
	}
	return NotFoundStock
}

// returns the sales stock of the given industry
// WAS 	err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", industry.Id, "Sales").StructScan(&stock)
func (industry Industry) SalesStock() Stock {
	username := industry.UserName
	stockList := (Users[username].StockList)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Sales`) {
			return *s
		}
	}
	return NotFoundStock
}

// returns the Labour Power stock of the given industry
// was query := `SELECT stocks.* FROM stocks INNER JOIN commodities ON stocks.commodity_id = commodities.id where stocks.owner_id = ? AND Usage_type ="Production" AND commodities.name="Labour Power"`
// bit of a botch to use the name of the commodity as a search term
func (industry Industry) VariableCapital() Stock {
	username := industry.UserName
	stockList := (Users[username].StockList)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Production`) && (s.CommodityName() == "Labour Power") {
			return *s
		}
	}
	return NotFoundStock
}

// returns the commodity that an industry produces
func (industry Industry) OutputCommodity() *Commodity {
	return industry.SalesStock().Commodity()
}

// return the productive capital stock of the given industry
// under development - at present assumes there is only one
// was 	query := `SELECT stocks.* FROM stocks INNER JOIN commodities ON stocks.commodity_id = commodities.id where stocks.owner_id = ? AND Usage_type ="Production" AND commodities.name="Means of Production"`
func (industry Industry) ConstantCapital() Stock {
	username := industry.UserName
	stockList := (Users[username].StockList)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Production`) && (s.CommodityName() == "Means of Production") {
			return *s
		}
	}
	return NotFoundStock
}

// returns all the constant capitals of a given industry
// TODO under development
// func (industry Industry) ConstantCapitals() []Stock {
// 	return &stocks [Programming error here]
// }

// METHODS OF SOCIAL CLASSES

// returns the sales stock of the given class
// was 	err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", class.Id, "Sales").StructScan(&stock)
func (class Class) MoneyStock() Stock {
	username := class.UserName
	stockList := (Users[username].StockList)

	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == class.Id) && (s.Usage_type == `Money`) {
			return *s
		}
	}
	return NotFoundStock
}

// returns the sales stock of the given class
func (class Class) SalesStock() Stock {
	username := class.UserName
	stockList := (Users[username].StockList)
	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == class.Id) && (s.Usage_type == `Sales`) {
			return *s
		}
	}
	return NotFoundStock
}

// returns the consumption stock of the given class
// under development - at present assumes there is only one
// WAS 	query := `SELECT stocks.* FROM stocks INNER JOIN commodities ON stocks.commodity_id = commodities.id where stocks.owner_id = ? AND Usage_type ="Consumption" AND commodities.name="Consumption"`
func (class Class) ConsumerGood() Stock {
	username := class.UserName
	stockList := (Users[username].StockList)

	for i := 0; i < len(stockList); i++ {
		s := &stockList[i]
		if (s.Owner_id == class.Id) && (s.Usage_type == `Consumption`) {
			return *s
		}
	}
	return NotFoundStock
}

// METHODS OF STOCKS

// returns a string specifying the stub of the URL that yields the owner of this stock
// if the stock is owned by an industry, this will link to an Industry object
// if the stock is owned by a social class, this URL will link to a class object
func (i Stock) OwnerLinkStub() string {
	switch i.Owner_type {
	case `Industry`:
		return `/industry`
	case `Class`:
		return `/class`
	default:
		return `unknown owner type`
	}
}

// fetches the name of the owner of this stock
func (s Stock) OwnerName() string {
	username := s.UserName
	switch s.Owner_type {
	case `Industry`:
		industryList := (Users[username].IndustryList)
		for i := 0; i < len(industryList); i++ {
			ind := &industryList[i]
			if s.Owner_id == ind.Id {
				return ind.Name
			}
		}
	case `Class`:
		classList := (Users[username].ClassList)
		for i := 0; i < len(classList); i++ {
			c := &classList[i]
			if s.Owner_id == c.Id {
				return c.Name
			}
		}
	default:
		return `UNKNOWN OWNER`
	}
	return `UNKNOWN OWNER`
}

// fetches the industry that owns this industry stock
// If it has none (an error, but we need to diagnose it) return nil.
func (s Industry_Stock) Industry() *Industry {
	industryList := (Users[s.UserName].IndustryList)
	for i := 0; i < len(industryList); i++ {
		ind := &industryList[i]
		if s.Industry_id == ind.Id {
			return ind
		}
	}
	return nil
}

// fetches the name of the industry that owns this industry stock.
// If it has none (an error, but we need to diagnose it) return "ERR"
func (s Industry_Stock) IndustryName() string {
	i := s.Industry()
	if i == nil {
		return "ERR"
	}
	return i.Name
}

// Return the name of the commodity that the given industry stock consists of.
// Return "UNKNOWN COMMODITY" if this is not found.
func (s Industry_Stock) CommodityName() string {
	username := s.UserName
	commodityList := (Users[username].CommodityList)
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return c.Name
		}
	}
	return `UNKNOWN COMMODITY`
}

// return the name of the commodity that the given stock consists of
// WAS 	rows, err := db.SDB.Queryx("SELECT * FROM commodities where Id = ?", i.Commodity_id)
func (s Stock) CommodityName() string {
	username := s.UserName
	commodityList := (Users[username].CommodityList)
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return c.Name
		}
	}
	return `UNKNOWN COMMODITY`
}

// return the commodity object that the given stock consists of
// WAS 	rows, err := db.SDB.Queryx("SELECT * FROM commodities where Id = ?", i.Commodity_id)
func (s Stock) Commodity() *Commodity {
	username := s.UserName
	commodityList := (Users[username].CommodityList)
	for i := 0; i < len(commodityList); i++ {
		c := commodityList[i]
		if s.Commodity_id == c.Id {
			return &c
		}
	}
	return &NotFoundCommodity
}

// under development
// will eventually be parameterised to yield value, price or quantity depending on a 'display' parameter
func (stock Stock) DisplaySize(mode string) float32 {
	switch mode {
	case `prices`:
		return stock.Size
	case `quantities`:
		return stock.Size // switch in price once this is in the model
	default:
		panic(`unknown display mode requested`)
	}
}

// (Experimental) Creates a url to link to this simulation, to be used in templates such as dashboard
// In this way all the URL naming is done in native Golang, not in the template
// We may also use such methods in the Trace function to improve usability
func (s Simulation) Link() string {
	return `/user/create/` + strconv.Itoa(s.Id)
}
