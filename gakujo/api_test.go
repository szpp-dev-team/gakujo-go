package gakujo

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var (
	username string
	password string
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("please set .env on ./..", err)
	}

	username = os.Getenv("J_USERNAME")
	password = os.Getenv("J_PASSWORD")
}

func TestLogin(t *testing.T) {
	c := NewClient()
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
}

func TestHome(t *testing.T) {
	c := NewClient()
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
	homeInfo, err := c.Home()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(homeInfo)
}
