package main

import (
	"bhbweb/api"
	"bhbweb/domain"
	"bhbweb/website"
	"fmt"
	"net/http"
)

func main() {
	book, err := domain.LoadBook("default")

	if err != nil {
		fmt.Println("Error reading book.")
		return
	}

	websiteController := website.NewController(book)
	apiAccountController := api.NewAccountController(book)
	apiTransactionController := api.NewTransactionController(book)

	http.HandleFunc("/home", websiteController.Home)
	http.HandleFunc("/addAccount", websiteController.AddAccount)
	http.HandleFunc("/addTransaction", websiteController.AddTransaction)
	http.HandleFunc("/viewAccount", websiteController.ViewAccount)

	http.HandleFunc("POST /api/account", apiAccountController.AddAccount)
	http.HandleFunc("POST /api/transaction", apiTransactionController.AddTransaction)

	http.ListenAndServe(":8080", nil)
}
