package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// The struct for the index page
type IndexData struct {
	Name    string
	Content BuisnessData
}

// Some business data I want to use in the template
// Might come from database or oher external service
type BuisnessData struct {
	NumRegisteredUsers int
	NumUsersToday      int
	TotalSales         int
	Profit             int
}

// fake data
var businessData = BuisnessData{NumRegisteredUsers: 5042, NumUsersToday: 48, TotalSales: 103452, Profit: 1932}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./index.html", "./global_components.html")

		if err != nil {
			fmt.Println(err.Error())
		}
		// render the index with the indexdata struct
		err = t.Execute(w, IndexData{"Alice", businessData})

		if err != nil {
			fmt.Println(err.Error())
		}

	})

	// The data endpoint to return only the business data
	// Might be for a seperate view - or for updating/swapping the component
	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./index.html", "./global_components.html")

		if err != nil {
			fmt.Println(err.Error())
		}
		// here i only the buisnessdata struct from the index page
		err = t.ExecuteTemplate(w, "business-data", businessData)

		if err != nil {
			fmt.Println(err.Error())
		}
	})

	// Endpoint for just seeing the profit
	r.HandleFunc("/profit", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./index.html", "./global_components.html")

		if err != nil {
			fmt.Println(err.Error())
		}
		// this template only needs the profit field (int)
		err = t.ExecuteTemplate(w, "profit", businessData.Profit)

		if err != nil {
			fmt.Println(err.Error())
		}
	})

	fmt.Println("Listening on 8080, open http://localhost:8080/")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}
