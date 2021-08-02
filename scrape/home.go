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

func NoticeDetail(r io.Reader) (model.NoticeDetail, error) {
	doc, _ := goquery.NewDocumentFromReader(r)
	txt := doc.Find("#right-box > form > div.right-module-bold.mt15 > div > div > table").Text()
	//txt = strings.Replace(txt, "カテゴリ", " ", -1)
	txt = strings.TrimSpace(txt)
	noticeDetail := parseDetialLine(txt)
	return noticeDetail, nil
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

func parseDetialLine(s string) model.NoticeDetail {
	count := 1
	element := []string{"カテゴリ", "タイトル", "連絡内容", "連絡元", "添付ファイル", "ファイルリンク公開", "参考URL", "重要度", "連絡日時", "WEB返信要求", "管理所属"}
	var noticedetail model.NoticeDetail
	for i := 0; i < len(element); i += 1 {
		idx := strings.Index(s, element[i])
		jdx := strings.Index(s, element[i+1])
		switch {
		case count == 1:
			noticedetail.Category = s[idx:jdx]

		case count == 2:
			noticedetail.Title = s[idx:jdx]

		case count == 3:
			noticedetail.Contact = s[idx:jdx]

		case count == 4:
			noticedetail.Detail = s[idx:jdx]

		case count == 5:
			noticedetail.Attachment = s[idx:strings.Index(s, "一括ダウンロード")]

		case count == 6:
			if Bool := strings.Index(s[idx:jdx-1], "公開する"); Bool > 0 {
				noticedetail.FilelinkPublication = true
			} else {
				noticedetail.FilelinkPublication = false
			}

		case count == 7:
			noticedetail.ReferenceURL = s[idx:jdx]

		case count == 8:
			if Bool := strings.Index(s[idx:jdx-1], "重要"); Bool > 0 {
				noticedetail.Important = true
			} else {
				noticedetail.Important = false
			}

		case count == 9:
			date, _ := time.Parse("2006/01/02", s[idx+len(element[i]):jdx-1])
			noticedetail.Date = date

		case count == 10:
			if Bool := strings.Index(s[idx:jdx-1], "返信を求めない"); Bool > 0 {
				noticedetail.WebReturnRequest = false
			} else {
				noticedetail.WebReturnRequest = true
			}

		case count == 11:
			noticedetail.Affiliation = s[idx:jdx]
		}
		count += 1
	}
	return noticedetail
}

/*func hoge() {
	strings.Replace(s string, "<div>", "", -1)
	strings.Replace(s string, "</div>", "\n", -1)
}*/
