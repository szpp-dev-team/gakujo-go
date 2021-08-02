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
	client   *Client
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("please set .env on ./..", err)
	}

	username = os.Getenv("J_USERNAME")
	password = os.Getenv("J_PASSWORD")
	begin = time.Now()
	client = NewClient()
}

func TestLogin(t *testing.T) {
	if err := client.Login(username, password); err != nil {
		t.Fatal(err)
	}
	t.Log("[Info]TestLogin Passed(took:", time.Since(begin), "ms)")
}

func TestLoadCookiesAndLogin(t *testing.T) {
	TestLogin(t)
	if err := client.DumpCookies(); err != nil {
		t.Fatal(err)
	}
	innerCli := NewClient()
	if err := innerCli.LoadCookiesAndLogin(); err != nil {
		t.Fatal(err)
	}
	t.Log("[Info]TestLoadCookiesAndLogin Passed(took:", time.Since(begin), "ms)")
}

func TestHome(t *testing.T) {
	TestLogin(t)

	homeInfo, err := client.Home()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(homeInfo)
	t.Log("[Info]TestHome Passed(took:", time.Since(begin), "ms)")
}
