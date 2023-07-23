package scrape

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-go/model"
	"github.com/szpp-dev-team/gakujo-go/util"
)

func ReportRows(r io.Reader) ([]model.ReportRow, error) {
	var err error
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	rows := []model.ReportRow{}
	selection := doc.Find("#searchList > tbody > tr")
	selection.EachWithBreak(func(i int, s *goquery.Selection) bool {
		courseName, courseDates, inerr := parseCourseNameFormat(s.Find("td:nth-child(1)").Text())
		if inerr != nil {
			err = inerr
			return false
		}

		title := util.ReplaceAndTrim(s.Find("td:nth-child(2)").Text())

		jsText, exists := s.Find("td:nth-child(2) > a").Attr("onclick")
		if !exists {
			err = errors.New("attr \"onClick\" not found")
			return false
		}
		reportMetadata, inerr := parseReportJSargument(jsText)
		if inerr != nil {
			err = inerr
			return false
		}

		statusText := util.ReplaceAndTrim(s.Find("td:nth-child(3)").Text())
		status := model.ToStatus(statusText)

		periodText := util.ReplaceAndTrim(s.Find("td:nth-child(4)").Text())
		beginDate, endDate, inerr := util.ParsePeriod(periodText)
		if inerr != nil {
			err = inerr
			return false
		}

		lastTimeText := util.ReplaceAndTrim(s.Find("td:nth-child(5)").Text())
		var lastSubmitDate time.Time
		if lastTimeText != "" {
			var lastTimeText1, lastTimeText2 string
			fmt.Sscanf(lastTimeText, "%s %s", &lastTimeText1, &lastTimeText2)
			lastSubmitDate, inerr = util.Parse2400("2006/01/02 15:04", fmt.Sprintf("%s %s", lastTimeText1, lastTimeText2))
			if inerr != nil {
				err = inerr
				return false
			}
		}

		format := util.ReplaceAndTrim(s.Find("td:nth-child(6)").Text())
		rows = append(rows, model.ReportRow{
			CourseName:     courseName,
			CourseDates:    courseDates,
			Title:          title,
			Status:         status,
			BeginDate:      beginDate,
			EndDate:        endDate,
			LastSubmitDate: lastSubmitDate,
			Format:         format,
			TaskMetadata:   reportMetadata,
		})
		return true
	})
	return rows, err
}

func ReportDetail(r io.Reader) (*model.ReportDetail, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	selection := doc.Find("#area > table > tbody")
	title := util.ReplaceAndTrim(selection.Find("tr:nth-child(1) > td").Text())
	periodText := util.ReplaceAndTrim(selection.Find("tr:nth-child(2) > td").Text())
	beginDate, endDate, err := util.ParsePeriod(periodText)
	if err != nil {
		return nil, err
	}
	evaluationMethod := util.ReplaceAndTrim(selection.Find("tr:nth-child(3) > td").Text())
	description, err := selection.Find("tr:nth-child(4) > td").Html()
	if err != nil {
		return nil, err
	}
	description = strings.Join(strings.Split(description, "<br/>"), "\n")
	transMatter := util.ReplaceAndTrim(selection.Find("tr:nth-child(6) > td").Text())
	return &model.ReportDetail{
		Title:            title,
		BeginDate:        beginDate,
		EndDate:          endDate,
		EvaluationMethod: evaluationMethod,
		Description:      description,
		TransMatter:      transMatter,
	}, nil
}

func parseReportJSargument(jsArgument string) (*model.TaskMetadata, error) {
	tokens := strings.Split(jsArgument[11:len(jsArgument)-2], ",")
	for i, token := range tokens {
		newToken := util.ReplaceAndTrim(token)
		tokens[i] = newToken[1 : len(newToken)-1]
	}
	if len(tokens) != 6 {
		return nil, errors.New("too few tokens")
	}

	year, err := strconv.Atoi(tokens[3])
	if err != nil {
		return nil, err
	}
	return &model.TaskMetadata{
		ID:               tokens[1],
		SubmitStatusCode: tokens[2],
		SchoolYear:       year,
		SubjectCode:      tokens[4],
		ClassCode:        tokens[5],
	}, nil
}
