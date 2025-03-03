package main

import (
	"log"
	"net/http"
	"os"

	"concurrent_money_transfer_system/internals/server"
	"concurrent_money_transfer_system/tests"
)

func main() {
	log.Println("Starting server on port 8080")
	router := server.SetupRouter()
	// Check if the with_test_users flag is provided
	withTestUsers := false
	for _, arg := range os.Args[1:] {
		if arg == "with_test_users" {
			withTestUsers = true
			break
		}
	}
	
	if withTestUsers {
		log.Println("Creating test users...")
		tests.CreateTestUsers()
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}