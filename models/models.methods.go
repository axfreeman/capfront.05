//methods.simulation.go
//class methods of the objects specified in models.simulation.go

package models

import "strconv"

// TODO eliminate boilerplate by making generic
// TODO see https://github.com/jose78/go-collection/blob/master/collections/collection.go  for suggestions

//METHODS OF INDUSTRIES

// A stock that will be returned if any condition is not met (that is, if the predicated stock does not exist)
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

// an aray of not found stocks
var NotFoundStocks = []Stock{NotFoundStock}

// returns the money stock of the given industry
// WAS err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", industry.Id, "Money").StructScan(&stock)

func (industry Industry) MoneyStock() Stock {
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Money`) {
			return *s
		}
	}
	return NotFoundStock
}

// returns the sales stock of the given industry
// WAS 	err = db.SDB.QueryRowx("SELECT * FROM stocks where Owner_Id = ? AND Usage_type =?", industry.Id, "Sales").StructScan(&stock)
func (industry Industry) SalesStock() Stock {
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
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
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
		if (s.Owner_id == industry.Id) && (s.Usage_type == `Production`) && (s.CommodityName() == "Labour Power") {
			return *s
		}
	}
	return NotFoundStock
}

// return the productive capital stock of the given industry
// under development - at present assumes there is only one
// was 	query := `SELECT stocks.* FROM stocks INNER JOIN commodities ON stocks.commodity_id = commodities.id where stocks.owner_id = ? AND Usage_type ="Production" AND commodities.name="Means of Production"`
func (industry Industry) ConstantCapital() Stock {
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
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
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
		if (s.Owner_id == class.Id) && (s.Usage_type == `Money`) {
			return *s
		}
	}
	return NotFoundStock
}

// returns the sales stock of the given class
func (class Class) SalesStock() Stock {
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
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
	for i := 0; i < len(StockList); i++ {
		s := &StockList[i]
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
// WAS err = db.SDB.QueryRowx("SELECT * FROM industries where Id = ?", stock.Owner_id).StructScan(&industry)
// WAS err = db.SDB.QueryRowx("SELECT * FROM classes where Id = ?", s.Owner_id).StructScan(&class)
func (s Stock) OwnerName() string {
	switch s.Owner_type {
	case `Industry`:
		for i := 0; i < len(IndustryList); i++ {
			ind := IndustryList[i]
			if s.Owner_id == ind.Id {
				return ind.Name
			}
		}
	case `Class`:
		for i := 0; i < len(ClassList); i++ {
			c := ClassList[i]
			if s.Owner_id == c.Id {
				return c.Name
			}
		}
	default:
		return `UNKNOWN OWNER`
	}
	return `UNKNOWN OWNER`
}

// return the name of the commodity that the given stock consists of
// WAS 	rows, err := db.SDB.Queryx("SELECT * FROM commodities where Id = ?", i.Commodity_id)
func (s Stock) CommodityName() string {
	for i := 0; i < len(CommodityList); i++ {
		c := CommodityList[i]
		if s.Commodity_id == c.Id {
			return c.Name
		}
	}
	return `UNKNOWN COMMODITY`
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
// We may use this technique for added clarity and perhaps to introduce into the Trace function to improve usability
func (s Simulation) Link() string {
	return `/user/select/` + strconv.Itoa(s.Id)
}
