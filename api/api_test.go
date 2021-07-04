package api

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLogin(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Fatal("please set .env on ./..", err)
	}

	username := os.Getenv("J_USERNAME")
	password := os.Getenv("J_PASSWORD")

	c := NewClient()
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
	if err := c.Home(); err != nil {
		t.Fatal(err)
	}
}
