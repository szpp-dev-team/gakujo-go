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
	c        *Client
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("please set .env on ./..", err)
	}

	username = os.Getenv("J_USERNAME")
	password = os.Getenv("J_PASSWORD")
	begin = time.Now()
	c = NewClient()
	if err := c.Login(username, password); err != nil {
		log.Fatal("failed to login")
	}
	log.Println("[Info]Login succeeded(took:", time.Since(begin), "ms)")
}

func TestLogin(t *testing.T) {
	inc := NewClient()
	if err := inc.Login(username, password); err != nil {
		t.Fatal(err)
	}
	t.Log("[Info]Login succeeded(took:", time.Since(begin), "ms)")
}

func TestHome(t *testing.T) {
	homeInfo, err := c.Home()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(homeInfo)
}

func TestSeisekiRows(t *testing.T) {
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

func TestDepartmentGpa(t *testing.T) {
	kc, err := c.NewKyoumuClient()
	if err != nil {
		t.Fatal(err)
	}
	dgpa, err := kc.DepartmentGpa()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*dgpa)
}
