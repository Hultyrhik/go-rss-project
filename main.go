package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hellow world")

	// loads .env file
	godotenv.Load(".env")

	// gets PORT value var
	portString := os.Getenv("PORT")
	if portString == "" {
		// exits program immediatly and logs a message
		log.Fatal("Port is not found in the enviroment")
	}

	fmt.Println("Port:", portString)
}
