package scrape

import (
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/util"
)

func ClassNoticeRow(r io.Reader) ([]model.ClassNoticeRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	replace := strings.NewReplacer("\t", "", "\n", "")
	classNoticeRows := make([]model.ClassNoticeRow, 0)
	doc.Find("#tbl_A01_01 > tbody > tr").Each(func(i int, selection *goquery.Selection) {

		targetDate := strings.TrimSpace(selection.Find("td:nth-child(6)").Text())

		dateText := selection.Find("td:nth-child(7)").Text()
		date, inerr := util.Parse2400("2006/01/02 15:04", dateText)
		if inerr != nil {
			err = inerr
			return
		}

		courseName := selection.Find("td:nth-child(2)").Text()
		teachersName := selection.Find("td:nth-child(3)").Text()
		title := selection.Find("td:nth-child(4)").Text()
		snt := selection.Find("td:nth-child(5)").Text()

		if len(targetDate) != 0 {
			targetdate, inerr := util.Parse2400("2006/01/02", targetDate)
			if inerr != nil {
				err = inerr
				return
			}
			classNoticeRow := model.ClassNoticeRow{
				CourseName:   strings.TrimSpace(replace.Replace(courseName)),
				TeachersName: strings.TrimSpace(teachersName),
				Title:        strings.TrimSpace(replace.Replace(title)),
				Type:         strings.TrimSpace(snt),
				TargetDate:   targetdate,
				Date:         date,
				Index:        i,
			}
			classNoticeRows = append(classNoticeRows, classNoticeRow)
		} else {
			var targetdate time.Time
			classNoticeRow := model.ClassNoticeRow{
				CourseName:   strings.TrimSpace(replace.Replace(courseName)),
				TeachersName: strings.TrimSpace(teachersName),
				Title:        strings.TrimSpace(replace.Replace(title)),
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
