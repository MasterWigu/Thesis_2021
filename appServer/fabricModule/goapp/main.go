package main

import (
	"log"
	"os"
)

func main() {

	createUsersBase()
	CreateProcessor()

	ginEngine := declareAPIEndpoint()

	log.Println("============ application-golang starts ============")
	err := os.Setenv("INITIALIZE-WITH-DISCOVERY", "false")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}


	startListening(ginEngine)

}
