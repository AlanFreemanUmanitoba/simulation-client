//models.objects.go
//describes the objects of the simulation itself
//functionally, these play two roles
// (1) they define how this front end communicates with the API of the backend
// (2) they define how this front end communicates with the user
// that is, the purpose is to intermediate between the simulation itself and the display of its results

package models

type Pair struct {
	Viewed   float32
	Compared float32
}

// This contains a record, generated by the server, of the results of the actions
type Trace struct {
	Id            int `json:"id" gorm:"primary_key"`
	Simulation_id int `json:"simulation_id"`
	Time_stamp    int
	UserName      string `json:"username"`
	Level         int    `json:"level"`
	Message       string `json:"message"`
}

// An individual simulation
// each user may have many simulations so there is a foreign key to the user
// each simulation contains a standard set of objects, being commodities, industries, classes and stocks
// so each object has a foreign key to the simulation
// all these keys are generated by the server and in the front end are used only to provide visual links in the display
// UserName is a convenience field which must be added by this frontend after downloading the data from the server
type Simulation struct {
	Id                     int    `json:"id"`
	Name                   string `json:"name"`
	Time_Stamp             int
	UserName               string  `json:"username"`
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
	Id                          int    `json:"id"`
	Name                        string `json:"name"`
	Simulation_id               int32  `json:"simulation_id"`
	Time_Stamp                  int32
	UserName                    string  `json:"username"`
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

type CommodityView struct {
	Id                          int
	Name                        string
	Origin                      string
	Usage                       string
	Size                        Pair
	Total_Value                 Pair
	Total_Price                 Pair
	Unit_Value                  Pair
	Unit_Price                  Pair
	Turnover_Time               Pair
	Demand                      Pair
	Supply                      Pair
	Allocation_Ratio            Pair
	Monetarily_Effective_Demand float32
	Investment_Proportion       float32
}

type Industry struct {
	Id                 int    `json:"id" gorm:"primary_key"`
	Name               string `json:"name"`
	Simulation_id      int32  `json:"simulation_id"`
	Time_Stamp         int
	UserName           string  `json:"username"`
	Output             string  `json:"output"`
	Output_Scale       float32 `json:"output_scale"`
	Output_Growth_Rate float32 `json:"output_growth_rate"`
	Initial_Capital    float32 `json:"initial_capital"`
	Work_In_Progress   float32 `json:"work_in_progress"`
	Current_Capital    float32 `json:"current_capital"`
	Profit             float32 `json:"profit"`
	Profit_Rate        float32 `json:"profit_rate"`
}

type IndustryView struct {
	Id                   int
	Name                 string
	OutputCommodityId    int
	Output               string
	Output_Scale         Pair
	Output_Growth_Rate   Pair
	Initial_Capital      Pair
	Work_In_Progress     Pair
	Current_Capital      Pair
	ConstantCapitalSize  Pair
	ConstantCapitalValue Pair
	ConstantCapitalPrice Pair
	VariableCapitalSize  Pair
	VariableCapitalValue Pair
	VariableCapitalPrice Pair
	MoneyStockSize       Pair
	MoneyStockValue      Pair
	MoneyStockPrice      Pair
	SalesStockSize       Pair
	SalesStockValue      Pair
	SalesStockPrice      Pair
	Profit               Pair
	Profit_Rate          Pair
}

type Class struct {
	Id                  int    `json:"id" gorm:"primary_key"`
	Name                string `json:"name"`
	Simulation_id       int32  `json:"simulation_id"`
	Time_Stamp          int
	UserName            string  `json:"username"`
	Population          float32 `json:"population"`
	Participation_Ratio float32 `json:"participation_ratio"`
	Consumption_Ratio   float32 `json:"consumption_ratio"`
	Revenue             float32 `json:"revenue"`
	Assets              float32 `json:"assets"`
}

type ClassView struct {
	Id                    int
	Name                  string
	Simulation_id         int32
	Time_Stamp            int
	UserName              string
	Population            Pair
	Participation_Ratio   float32
	Consumption_Ratio     float32
	Revenue               Pair
	Assets                Pair
	ConsumptionStockSize  Pair
	ConsumptionStockValue Pair
	ConsumptionStockPrice Pair
	MoneyStockSize        Pair
	MoneyStockValue       Pair
	MoneyStockPrice       Pair
	SalesStockSize        Pair
	SalesStockValue       Pair
	SalesStockPrice       Pair
}

type Industry_Stock struct {
	Id            int     `json:"id" gorm:"primary_key"`
	Simulation_id int     `json:"simulation_id" `
	Industry_id   int     `json:"industry_id"`
	Commodity_id  int     `json:"commodity_id" `
	UserName      string  `json:"username"`
	Name          string  `json:"name" `
	Usage_type    string  `json:"usage_type" `
	Size          float32 `json:"size" `
	Value         float32 `json:"value" `
	Price         float32 `json:"price" `
	Requirement   float32 `json:"requirement" `
	Demand        float32 `json:"demand" `
}

type Class_Stock struct {
	Id            int     `json:"id" gorm:"primary_key"`
	Simulation_id int     `json:"simulation_id" `
	Class_id      int     `json:"class_id"`
	Commodity_id  int     `json:"commodity_id"`
	UserName      string  `json:"username"`
	Name          string  `json:"name" `
	Usage_type    string  `json:"usage_type" `
	Size          float32 `json:"size" `
	Value         float32 `json:"value" `
	Price         float32 `json:"price" `
	Demand        float32 `json:"demand" `
}

// This list of templates is common to all users.
// It would normally change only when the database is reset from
// immutable fixtures using Refresh().
// It is initialized when this frontend restarts.
// In future there should be some procedure for adding new templates
// or editing existing ones.
var TemplateList []Simulation

// a HistoryItem contains all the information describing a stage
// of the Simulation. The UserData object contains a map[int]HistoryItem
// which lets the user review past stages of the current Simulation
type HistoryItem struct {
	SimulationList    []Simulation
	CommodityList     []Commodity
	IndustryList      []Industry
	ClassList         []Class
	IndustryStockList []Industry_Stock
	ClassStockList    []Class_Stock
	TraceList         []Trace
	Time_stamp        int
	State             string
}

func NewHistoryItem() HistoryItem {
	NewItem := HistoryItem{
		Time_stamp:        0,
		State:             "Demand",
		SimulationList:    make([]Simulation, 0),
		CommodityList:     make([]Commodity, 0),
		IndustryList:      make([]Industry, 0),
		ClassList:         make([]Class, 0),
		IndustryStockList: make([]Industry_Stock, 0),
		ClassStockList:    make([]Class_Stock, 0),
		TraceList:         make([]Trace, 0),
	}
	return NewItem
}

var SimulationList []Simulation
