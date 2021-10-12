package gakujo

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/util"
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

func TestClassNoticeRow(t *testing.T) {
	opt := model.BasicClassNoticeSearchOpt(2021, model.ToSemesterCode("前期"), util.BasicTime(2021, 3, 1))
	classNoticeRow, err := c.ClassNoticeRows(opt)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range classNoticeRow {
		fmt.Printf("%+v\n", row)
	}
}

func TestAllClassNoticeRow(t *testing.T) {
	opt := model.AllClassNoticeSearchOpt(2020)
	classNoticeRow, err := c.ClassNoticeRows(opt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v records found\n", len(classNoticeRow))
	for _, row := range classNoticeRow {
		fmt.Printf("%+v\n", row)
	}
}

func TestClassNoticeDetail(t *testing.T) {
	opt := model.BasicClassNoticeSearchOpt(2021, model.ToSemesterCode("前期"), util.BasicTime(2021, 3, 1))
	classNoticeRow, err := c.ClassNoticeRows(opt)
	if err != nil {
		t.Fatal(err)
	}
	detail, err := c.ClassNoticeDetail(&classNoticeRow[0], opt)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(*detail)
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

func TestChusenRegistrationRows(t *testing.T) {
	kc, err := c.NewKyoumuClient()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := kc.ChusenRegistrationRows()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		fmt.Printf("%+v\n", *row)
	}
}

func TestPostChusenRegistration(t *testing.T) {
	kc, err := c.NewKyoumuClient()
	if err != nil {
		t.Fatal(err)
	}

	rows, err := kc.ChusenRegistrationRows()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("=========== original ===========")
	fmt.Printf("科目名: %s 第%d志望\n\n", rows[0].SubjectName, rows[0].ChoiceRank)

	rows[0].ChoiceRank = 0
	if err := kc.PostChusenRegistration(rows); err != nil {
		t.Fatal(err)
	}

	rows, err = kc.ChusenRegistrationRows()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("=========== new ===========")
	fmt.Printf("科目名: %s 第%d志望\n\n", rows[0].SubjectName, rows[0].ChoiceRank)
}

func TestPostRishuRegistration(t *testing.T) {
	// 火曜日・2限のとある科目(CSB2限定)
	const (
		kamokuCode = "77605250"
		classCode  = "61"
		unit       = 2
		radio      = 0
		youbi      = 2
		jigen      = 2
	)
	kc, err := c.NewKyoumuClient()
	if err != nil {
		t.Fatal(err)
	}

	formData := model.NewPostKamokuFormData(kamokuCode, classCode, unit, radio, youbi, jigen)
	if err := kc.PostRishuRegistration(formData); err != nil {
		t.Fatal(err)
	}
}

func TestReportRows(t *testing.T) {
	option := model.ReportSearchOption{
		SchoolYear:   2020,
		SemesterCode: model.EarlyPeriod,
	}
	reportRows, err := c.ReportRows(&option)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range reportRows {
		fmt.Println(row)
	}
}

func TestReportDetail(t *testing.T) {
	option := model.ReportSearchOption{
		SchoolYear:   2020,
		SemesterCode: model.EarlyPeriod,
	}
	rows, err := c.ReportRows(&option)
	if err != nil {
		t.Fatal(err)
	}
	detailOption := rows[0].DetailOption()
	reportDetail, err := c.ReportDetail(detailOption)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(reportDetail)
}

func TestMinitestRows(t *testing.T) {
	option := model.MinitestSearchOption{
		SchoolYear:   2020,
		SemesterCode: model.EarlyPeriod,
	}
	rows, err := c.MinitestRows(&option)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		fmt.Println(row)
	}
}
