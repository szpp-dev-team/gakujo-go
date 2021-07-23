package scrape

import (
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
)

func TaskRows(r io.Reader) ([]model.TaskRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	taskRows := make([]model.TaskRow, 0)
	doc.Find("#tbl_submission > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		taskType, inerr := model.ToTasktype(selection.Find("td.arart > span > span").Text())
		if inerr != nil {
			err = inerr
			return
		}
		deadlineText := selection.Find("td.daytime").Text()
		var deadline time.Time
		if deadlineText != "" {
			deadline, inerr = time.Parse("2006/01/02 15:04", deadlineText)
			if inerr != nil {
				err = inerr
				return
			}
		}
		taskRow := model.TaskRow{
			Type:     taskType,
			Deadline: deadline,
			Name:     selection.Find("td:nth-child(3) > a").Text(),
			Index:    i,
		}
		taskRows = append(taskRows, taskRow)
	})
	if err != nil {
		return nil, err
	}
	return taskRows, nil
}

func NoticeRows(r io.Reader) ([]model.NoticeRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	noticeRows := make([]model.NoticeRow, 0)
	doc.Find("#tbl_news > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		noticeType, inerr := model.ToNoticetype(selection.Find("td.arart > span > span > a").Text())
		if inerr != nil {
			err = inerr
			return
		}
		titleLine := selection.Find("td.title > a").Text()
		snt, important, title, inerr := parseTitleLine(titleLine)
		if inerr != nil {
			err = inerr
			return
		}
		dateText := selection.Find("td.day").Text()
		date, inerr := time.Parse("2006/01/02", dateText)
		if inerr != nil {
			err = inerr
			return
		}
		noticeRow := model.NoticeRow{
			Type:        noticeType,
			SubType:     snt,
			Important:   important,
			Title:       title,
			Date:        date,
			Affiliation: selection.Find("td.syozoku").Text(),
			Index:       i,
		}
		noticeRows = append(noticeRows, noticeRow)
	})
	if err != nil {
		return nil, err
	}
	return noticeRows, nil
}

// return (SubNoticeType, isImportant, title)
func parseTitleLine(s string) (model.SubNoticeType, bool, string, error) {
	big := false
	squ := false
	bigText := ""
	squText := ""
	important := false
	title := ""
	for _, c := range s {
		if c == '【' {
			big = true
			continue
		}
		if c == '】' {
			big = false
			continue
		}
		if c == '[' {
			squ = true
			continue
		}
		if c == ']' {
			squ = false
			continue
		}
		if big {
			bigText += string(c)
		} else if squ {
			squText += string(c)
		} else {
			title += string(c)
		}
	}
	if bigText == "重要" {
		important = true
	}
	return model.ToSubNoticetype(squText), important, title, nil
}