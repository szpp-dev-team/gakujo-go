package gakujo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/szpp-dev-team/gakujo-api/scrape"
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

func TestNoticeDetail(t *testing.T) {
	c := NewClient()
	if err := c.Login(username, password); err != nil {
		t.Fatal(err)
	}
	t.Log("[Info]Login succeeded(took:", time.Since(begin), "ms)")
	noticeDetailHtml, _ := c.fetchNoiceDetailhtml()
	defer noticeDetailHtml.Close()
	txt, _ := scrape.NoticeDetail(noticeDetailHtml)
	fmt.Println(txt)
}
