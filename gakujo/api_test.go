package gakujo

import (
	"fmt"
	"log"
	"os"
	"strings"
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
	noticeDetail, err := scrape.NoticeDetail(noticeDetailHtml)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 11; i++ {
		switch {
		case i == 0:
			fmt.Print("カテゴリー  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Category))
		case i == 1:
			fmt.Print("タイトル  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Title))
		case i == 2:
			fmt.Print("連絡内容  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Detail))
		case i == 3:
			fmt.Print("連絡元  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Contact))
		case i == 4:
			fmt.Print("添付ファイル  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Attachment))
		case i == 5:
			fmt.Print("ファイルリンク公開  ")
			fmt.Println(noticeDetail.FilelinkPublication)
		case i == 6:
			fmt.Print("参照URL  ")
			fmt.Println(strings.TrimSpace(noticeDetail.ReferenceURL))
		case i == 7:
			fmt.Print("重要度  ")
			fmt.Println(noticeDetail.Important)
		case i == 8:
			fmt.Print("日時  ")
			fmt.Println(noticeDetail.Date)
		case i == 9:
			fmt.Print("WEB返信要求  ")
			fmt.Println(noticeDetail.WebReturnRequest)
		case i == 10:
			fmt.Print("管理所属  ")
			fmt.Println(strings.TrimSpace(noticeDetail.Affiliation))
		}

	}
}
