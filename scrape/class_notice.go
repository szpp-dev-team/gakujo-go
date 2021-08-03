package scrape

import (
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
)

func ClassNoticeRow(r io.Reader) ([]model.ClassNoticeRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	classNoticeRows := make([]model.ClassNoticeRow, 0)
	doc.Find("#tbl_A01_01 > tbody > tr").Each(func(i int, selection *goquery.Selection) {

		targetDate := strings.TrimSpace(selection.Find("td:nth-child(6)").Text())

		dateText := selection.Find("td:nth-child(7)").Text()
		date, inerr := time.Parse("2006/01/02 15:04", dateText)
		if inerr != nil {
			err = inerr
			return
		}

		courses := selection.Find("td:nth-child(2)").Text()
		teachersName := selection.Find("td:nth-child(3)").Text()
		title := selection.Find("td:nth-child(4)").Text()
		snt := selection.Find("td:nth-child(5)").Text()

		if len(targetDate) != 0 {
			targetdate, inerr := time.Parse("2006/01/02", targetDate)
			if inerr != nil {
				err = inerr
				return
			}
			classNoticeRow := model.ClassNoticeRow{
				Courses:      strings.TrimSpace(courses),
				TeachersName: strings.TrimSpace(teachersName),
				Title:        strings.TrimSpace(title),
				Type:         strings.TrimSpace(snt),
				TargetDate:   targetdate,
				Date:         date,
				Index:        i,
			}
			classNoticeRows = append(classNoticeRows, classNoticeRow)
		} else {
			dateText := selection.Find("td:nth-child(7)").Text()
			date, inerr := time.Parse("2006/01/02 15:04", dateText)
			if inerr != nil {
				err = inerr
				return
			}
			var targetdate time.Time
			classNoticeRow := model.ClassNoticeRow{
				Courses:      strings.TrimSpace(courses),
				TeachersName: strings.TrimSpace(teachersName),
				Title:        strings.TrimSpace(title),
				Type:         strings.TrimSpace(snt),
				TargetDate:   targetdate,
				Date:         date,
				Index:        i,
			}
			classNoticeRows = append(classNoticeRows, classNoticeRow)
		}
	})
	if err != nil {
		return nil, err
	}
	return classNoticeRows, err
}

func ClassNoticeRow2(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	return strings.TrimSpace(doc.Find("#tbl_A01_01 > tbody > tr:nth-child(1) > td:nth-child(6)").Text()), err
}
