package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/szpp-dev-team/gakujo-api/api"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("please set .env on ./..", err)
	}

	username := os.Getenv("J_USERNAME")
	password := os.Getenv("J_PASSWORD")

	fmt.Println(username, password)

	c := api.NewClient()
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	if err := c.Home(); err != nil {
		log.Fatal(err)
	}
}
