package scrape

import (
	"io"
	"strings"
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

func NoticeDetail(r io.Reader) (*model.NoticeDetail, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var noticeDetail model.NoticeDetail
	doc.Find("#right-box > form > div.right-module-bold.mt15 > div > div > div > table.ttb_entry > tbody > tr").Each(func(i int, selection *goquery.Selection) {
		switch {
		case i == 0:
			noticeDetail.ContactType = strings.TrimSpace(selection.Find("td").Text())
		case i == 1:
			noticeDetail.Title = strings.TrimSpace(selection.Find("td").Text())
		case i == 2:
			noticeDetail.Detail = strings.TrimSpace(selection.Find("td").Text())
		case i == 3:
			noticeDetail.File = strings.TrimSpace(selection.Find("td").Text())
		case i == 4:
			noticeDetail.FilelinkPublication = !strings.Contains(selection.Find("td").Text(), "公開しない")
		case i == 5:
			noticeDetail.ReferenceURL = strings.TrimSpace(selection.Find("td").Text())
		case i == 6:
			noticeDetail.Important = !strings.Contains(selection.Find("td").Text(), "通常")
		case i == 7:
			rawText := strings.Replace(selection.Find("td").Text(), "即時通知", "", -1)
			rawText = strings.TrimSpace(rawText)
			date, inerr := time.Parse("2006/01/02 15:04", rawText)
			if inerr != nil {
				err = inerr
				return
			}
			noticeDetail.Date = date
		case i == 8:
			noticeDetail.WebReturnRequest = !strings.Contains(selection.Find("td").Text(), "返信を求めない")
		}
	})
	if err != nil {
		return nil, err
	}

	return &noticeDetail, nil
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
