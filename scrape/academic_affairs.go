package scrape

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-api/model"
)

func SeisekiRows(r io.Reader) ([]*model.SeisekiRow, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	seisekiRows := make([]*model.SeisekiRow, 0)
	doc.Find("table.txt12 > tbody:nth-child(1) > tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i == 0 {
			return true
		}
		seisekiRow, inerr := scrapeSeisekiTrRow(s)
		if inerr != nil {
			err = inerr
			return false
		}
		seisekiRows = append(seisekiRows, seisekiRow)
		return true
	})
	return seisekiRows, err
}

func scrapeSeisekiTrRow(s *goquery.Selection) (*model.SeisekiRow, error) {
	var (
		seisekiRow model.SeisekiRow
		err        error
	)
	replacer := strings.NewReplacer("\n", "", "\t", "", " ", "")

	s.Find("td").EachWithBreak(func(i int, ins *goquery.Selection) bool {
		if seisekiRow.Grade == model.GradeType(model.GTPassed) && 6 <= i && i <= 8 {
			return true
		}
		elm := func() *goquery.Selection {
			if i >= 4 {
				return ins
			} else {
				return ins.Find("span:nth-child(1)")
			}
		}()

		text := elm.Text()
		value := replacer.Replace(elm.AttrOr("title", text))
		switch i {
		case 0:
			seisekiRow.SubjectName = value
		case 1:
			seisekiRow.TeacherName = value
		case 2:
			seisekiRow.SubjectDistinction = value
		case 3:
			st, inerr := model.ToSubjectType(value)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.SubjectType = st
		case 4:
			credit, inerr := strconv.Atoi(value)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.Credit = credit
		case 5:
			gt, inerr := model.ToGradeType(value)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.Grade = gt
		case 6:
			score, inerr := strconv.ParseFloat(value, 64)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.Score = score
		case 7:
			gp, inerr := strconv.ParseFloat(value, 64)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.GP = gp
		case 8:
			year, inerr := strconv.Atoi(value)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.Year = year
		case 9:
			date, inerr := time.Parse("2006-01-02", value)
			if inerr != nil {
				err = inerr
				return false
			}
			seisekiRow.Date = date
		default:
			fmt.Printf("i: %d was ignored\n", i)
		}
		return true
	})

	return &seisekiRow, err
}
