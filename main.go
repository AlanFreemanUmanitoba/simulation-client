package main

import (
	"capfront/display"
	"capfront/fetch"
	"fmt"
	"net/http"
)

func main() {

	display.Router.LoadHTMLGlob("./templates/**/*") // load the templates

	fmt.Println("Welcome to capitalism")

	dir := http.Dir("./static")
	fs := http.FileServer(dir)
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	display.Router.GET("/action/:action", display.ActionHandler)

	display.Router.GET("/commodities", display.ShowCommodities)
	display.Router.GET("/industries", display.ShowIndustries)
	display.Router.GET("/classes", display.ShowClasses)
	display.Router.GET("/industry_stocks", display.ShowIndustryStocks)
	display.Router.GET("/class_stocks", display.ShowClassStocks)
	display.Router.GET("/trace", display.ShowTrace)

	display.Router.GET("/industry/:id", display.ShowIndustry)
	display.Router.GET("/commodity/:id", display.ShowCommodity)
	display.Router.GET("/class/:id", display.ShowClass)

	display.Router.GET("/user/create/:id", display.CreateSimulation)
	display.Router.GET("/user/switch/:id", display.SwitchSimulation)
	display.Router.GET("/user/delete/:id", display.DeleteSimulation)
	display.Router.GET("/user/restart/:id", display.RestartSimulation)

	display.Router.GET("/", display.ShowIndexPage)
	display.Router.GET("/data/", display.DataHandler)
	display.Router.GET("/user/dashboard", display.UserDashboard)
	display.Router.GET("/admin/dashboard", display.AdminDashboard)
	display.Router.GET("/admin/reset", display.AdminReset)
	display.Router.GET("/admin/choose-players", display.Lock)

	display.Router.GET("/back", display.Back)
	display.Router.GET("/forward", display.Forward)

	fetch.Initialise()
	display.ListData()

	display.Router.Run() // Run the server

}
