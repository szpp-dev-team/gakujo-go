package scrape

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/util"
)

func ClassEnqRows(r io.Reader) ([]model.ClassEnqRow, error) {
	var err error
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	rows := []model.ClassEnqRow{}
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
			err = errors.New("Attr \"onClick\" not found")
			return false
		}
		ClassEnqMetadata, inerr := parseClassEnqJSargument(jsText)
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

		submitStatusText := util.ReplaceAndTrim(s.Find("td:nth-child(5)").Text())
		submitStatus := model.ToSubmitStatus(submitStatusText)

		rows = append(rows, model.ClassEnqRow{
			CourseName:       courseName,
			CourseDates:      courseDates,
			Title:            title,
			Status:           status,
			BeginDate:        beginDate,
			EndDate:          endDate,
			SubmitStatus:     submitStatus,
			ClassEnqMetadata: ClassEnqMetadata,
		})
		return true
	})
	return rows, err
}

func ClassEnqDetail(r io.Reader) (model.ClassEnqDetail, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return model.ClassEnqDetail{}, err
	}
	selection := doc.Find("#right-box > form:nth-child(4) > div.right-module-bold.mt15 > div > div > table:nth-child(1) > tbody")
	title := util.ReplaceAndTrim(selection.Find("tr:nth-child(1) > td").Text())
	periodText := util.ReplaceAndTrim(selection.Find("tr:nth-child(2) > td").Text())
	beginDate, endDate, err := util.ParsePeriod(periodText)
	if err != nil {
		return model.ClassEnqDetail{}, err
	}
	numText := util.ReplaceAndTrim(selection.Find("tr:nth-child(3) > td").Text())
	var num int
	fmt.Sscanf(numText, "%d 問", &num)
	description := util.ReplaceAndTrim(selection.Find("tr:nth-child(5) > td").Text())
	description = strings.Join(strings.Split(description, "<br/>"), "\n")
	transMatter := util.ReplaceAndTrim(selection.Find("tr:nth-child(7) > td").Text())
	ClassEnqHtml, err := doc.Find("#area > div:nth-child(4)").Html() // not working
	if err != nil {
		return model.ClassEnqDetail{}, err
	}
	return model.ClassEnqDetail{
		Title:        title,
		BeginDate:    beginDate,
		EndDate:      endDate,
		Num:          num,
		Description:  description,
		TransMatter:  transMatter,
		ClassEnqHtml: ClassEnqHtml,
	}, nil
}

func parseClassEnqJSargument(jsArgument string) (model.ClassEnqMetadata, error) {
	tokens := strings.Split(jsArgument[11:len(jsArgument)-2], ",")
	for i, token := range tokens {
		newToken := util.ReplaceAndTrim(token)
		tokens[i] = newToken[1 : len(newToken)-1]
	}
	if len(tokens) != 6 {
		return model.ClassEnqMetadata{}, errors.New("Too few tokens")
	}

	year, err := strconv.Atoi(tokens[3])
	if err != nil {
		return model.ClassEnqMetadata{}, err
	}
	return model.ClassEnqMetadata{
		ClassEnqID:       tokens[1],
		SubmitStatusCode: tokens[2],
		SchoolYear:       year,
		SubjectCode:      tokens[4],
		ClassCode:        tokens[5],
	}, nil
}
