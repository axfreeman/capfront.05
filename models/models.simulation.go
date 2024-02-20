//models.simulation.go
//describes the objects of the simulation itself
//functionally, these play two roles
// (1) they define how this front end communicates with the API of the backend
// (2) they define how this front end communicates with the user
// that is, the purpose is to intermediate between the simulation itself and the display of its results

package models

type Trace struct {
	Id            int `json:"id" gorm:"primary_key"`
	Simulation_id int `json:"simulation_id"`
	Time_stamp    int
	Level         int    `json:"level"`
	Message       string `json:"message"`
}

type Simulation struct {
	Id                     int    `json:"id" gorm:"primary_key"`
	Name                   string `json:"name"`
	Time_Stamp             int
	State                  string  `json:"state"`
	Periods_Per_Year       float32 `json:"periods_per_year"`
	Population_Growth_Rate float32 `json:"population_growth_rate"`
	Investment_Ratio       float32 `json:"investment_ratio"`
	Labour_Supply_Demand   string  `json:"labour_supply_response"`
	Price_Response_Type    string  `json:"price_response_type"`
	Melt_Response_Type     string  `json:"melt_response_type"`
	Currency_Symbol        string  `json:"currency_symbol"`
	Quantity_Symbol        string  `json:"quantity_symbol"`
	Melt                   float32 `json:"melt"`
	User                   int32   `json:"user_id"`
}

type Commodity struct {
	Id                          int    `json:"id" gorm:"primary_key"`
	Name                        string `json:"name"`
	Simulation                  int32  `json:"simulation_id"`
	Time_Stamp                  int32
	Origin                      string  `json:"origin"`
	Usage                       string  `json:"usage"`
	Size                        float32 `json:"size"`
	Total_Value                 float32 `json:"total_value"`
	Total_Price                 float32 `json:"total_price"`
	Unit_Value                  float32 `json:"unit_value"`
	Unit_Price                  float32 `json:"unit_price"`
	Turnover_Time               float32 `json:"turnover_time"`
	Demand                      float32 `json:"demand"`
	Supply                      float32 `json:"supply"`
	Allocation_Ratio            float32 `json:"allocation_ratio"`
	Display_Order               float32 `json:"display_order"`
	Image_Name                  string  `json:"image_name"`
	Tooltip                     string  `json:"tooltip"`
	Monetarily_Effective_Demand float32 `json:"monetarily_effective_demand"`
	Investment_Proportion       float32 `json:"investment_proportion"`
}

type Industry struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	Name               string `json:"name"`
	Simulation         int32  `json:"simulation_id"`
	Time_Stamp         int
	Output             string  `json:"output"`
	Output_Scale       float32 `json:"output_scale"`
	Output_Growth_Rate float32 `json:"output_growth_rate"`
	Initial_Capital    float32 `json:"initial_capital"`
	Work_In_Progress   float32 `json:"work_in_progress"`
	Current_Capital    float32 `json:"current_capital"`
	Profit             float32 `json:"profit"`
	Profit_Rate        float32 `json:"profit_rate"`
}

type Class struct {
	Id                  int    `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	Simulation          int32  `json:"simulation_id"`
	Time_Stamp          int
	Population          float32 `json:"population"`
	Participation_Ratio float32 `json:"participation_ratio"`
	Consumption_Ratio   float32 `json:"consumption_ratio"`
	Revenue             float32 `json:"revenue"`
	Assets              float32 `json:"assets"`
}

type Stock struct {
	Id            int `json:"id" gorm:"primary_key"`
	Simulation_id int `json:"simulation_id" ` //TODO consistently rename all other FKs to Simulation_id
	Time_Stamp    int
	Owner_id      int     `json:"owner_id"`
	Commodity_id  int     `json:"commodity_id" `
	Name          string  `json:"name" `
	Owner_type    string  `json:"owner_type" `
	Usage_type    string  `json:"usage_type" `
	Size          float32 `json:"size" `
	Value         float32 `json:"value" `
	Price         float32 `json:"price" `
	Requirement   float32 `json:"requirement" `
	Demand        float32 `json:"demand" `
}

var SimulationList []Simulation
var CommodityList []Commodity
var IndustryList []Industry
var ClassList []Class
var StockList []Stock
var TraceList []Trace
var MySimulations []Simulation //TODO prototype of in-memory handling which should simplify matters if it does not slow things down too much
