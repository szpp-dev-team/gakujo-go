package scrape

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
)

func ReportRows(r io.Reader) ([]model.ReportRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	reportRows := make([]model.ReportRow, 0, 32)

	var scrapeError error = nil

	doc.Find("#searchList > tbody > tr").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		subject := sel.Find("td:nth-child(1)").Text()
		title := sel.Find("td:nth-child(2)").Text()
		stateText := sel.Find("td:nth-child(3)").Text()
		submissionPeriodText := sel.Find("td:nth-child(4)").Text()
		lastSubmissionTimeText := sel.Find("td:nth-child(5)").Text()
		formatText := sel.Find("td:nth-child(6)").Text()

		state, _ := model.ToReportState(strings.TrimSpace(stateText))
		format, _ := model.ToReportFormat(strings.TrimSpace(formatText))

		startTime, dueTime, err := parseReportSubmissionPeriod(submissionPeriodText)
		if err != nil {
			scrapeError = err
			return false
		}

		lastSubmissionTime, err := parseLastSubmissionTime(lastSubmissionTimeText)
		if err != nil {
			scrapeError = err
			return false
		}

		row := model.ReportRow{
			Subject:            strings.TrimSpace(subject),
			Title:              strings.TrimSpace(title),
			State:              state,
			StartTime:          startTime,
			DueTime:            dueTime,
			LastSubmissionTime: lastSubmissionTime,
			Format:             format,
		}
		reportRows = append(reportRows, row)

		return true
	})

	if scrapeError != nil {
		return nil, scrapeError
	}
	return reportRows, nil
}

func parseReportSubmissionPeriod(s string) (start, due time.Time, err error) {
	const layout = "2006/01/02 15:04"
	loc, _ := time.LoadLocation("Asia/Tokyo")

	parts := strings.Split(s, " ")
	if len(parts) != 5 {
		return time.Time{}, time.Time{},
			fmt.Errorf("%q is an invalid format for report submission period", s)
	}

	hour24 := strings.HasPrefix(parts[1], "24")
	if hour24 {
		parts[1] = "00" + parts[1][2:]
	}
	start, err = time.ParseInLocation(layout, parts[0]+" "+parts[1], loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if hour24 {
		start = start.Add(time.Duration(24) * time.Hour)
	}

	hour24 = strings.HasPrefix(parts[4], "24")
	if hour24 {
		parts[4] = "00" + parts[4][2:]
	}
	due, err = time.ParseInLocation(layout, parts[3]+" "+parts[4], loc)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if hour24 {
		due = due.Add(time.Duration(24) * time.Hour)
	}

	return start, due, nil
}

func parseLastSubmissionTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return time.Time{}, nil
	}

	return time.Parse("2006/01/02 15:04", s)
}
