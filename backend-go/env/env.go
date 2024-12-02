package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var BE_PORT string
var VALKEY_ADDRESS string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	BE_PORT = os.Getenv("BE_PORT")
	VALKEY_ADDRESS = os.Getenv("VALKEY_ADDRESS")
}
