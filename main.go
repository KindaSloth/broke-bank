package main

import (
	"broke-bank/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	addr, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		log.Fatal("Missing SERVER_ADDRESS env")
	}

	s := server.New()
	s.Run(addr)
}
