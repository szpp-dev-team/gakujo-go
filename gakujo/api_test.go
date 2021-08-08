package gakujo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

var (
	begin    time.Time
	username string
	password string
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("please set .env on ./..", err)
	}

	username = os.Getenv("J_USERNAME")
	password = os.Getenv("J_PASSWORD")
	begin = time.Now()
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
	t.Log("[Info]Login succeeded(took:", time.Since(begin), "ms)")
	homeInfo, err := c.Home()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(homeInfo)
}

func TestSeisekiRows(t *testing.T) {
	c := NewClient()
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
	t.Log("[Info]Login succeeded(took:", time.Since(begin), "ms)")

	kc, err := c.NewKyoumuClient()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := kc.SeisekiRows()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		fmt.Println(*row)
	}
}
