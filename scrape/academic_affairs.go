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

func DepartmentGpa(r io.Reader) (*model.DepartmentGpa, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return scrapeGpaTrRow(doc.Find("table.txt12 > tbody:nth-child(1)"))
}

func scrapeGpaTrRow(s *goquery.Selection) (*model.DepartmentGpa, error) {
	var (
		departmentGpa model.DepartmentGpa
		termGPA       model.TermGpa = model.TermGpa{}
		err           error
	)
	replacer := strings.NewReplacer("\n", "", "\t", "", " ", "")
	s.Find("tr").EachWithBreak(func(i int, ins *goquery.Selection) bool {
		if i == 0 {
			return true
		}
		item := ins.Find("td:nth-child(1) > font").Text()
		text := replacer.Replace(ins.Find("td:nth-child(2)").Text())

		if strings.Contains(item, "GPA値") {
			var (
				year int
				term string
			)
			gpa, _ := strconv.ParseFloat(text, 64)
			if _, err := fmt.Sscanf(item, "%d年度　%s　GPA値", &year, &term); err != nil {
				departmentGpa.Gpa = gpa
			} else {
				termGPA.Year = year
				if term == "前期" {
					termGPA.FirstGPA = gpa
				} else if term == "後期" {
					termGPA.SecondGPA = gpa

					departmentGpa.TermGpas = append(departmentGpa.TermGpas, termGPA)
					termGPA = model.TermGpa{}
				}
			}
		} else {
			if termGPA.FirstGPA != 0 {
				departmentGpa.TermGpas = append(departmentGpa.TermGpas, termGPA)
				termGPA = model.TermGpa{}
			}

			switch item {
			case "最終GPA算出日":
				t, inerr := time.Parse("2006年01月02日", text)
				if inerr != nil {
					err = inerr
					return false
				}
				departmentGpa.CalcDate = t
			case "同一学科内順位", "同一コース内順位":
				var rank, num int
				if _, inerr := fmt.Sscanf(text, "%d人中　%d位", &num, &rank); inerr != nil {
					err = inerr
					return false
				}
				if strings.Contains(item, "学科") {
					departmentGpa.DepartmentNum = num
					departmentGpa.DepartmentRank = rank
				} else {
					departmentGpa.CourseNum = num
					departmentGpa.CourseRank = rank
				}
			}
		}

		return true
	})

	return &departmentGpa, err
}

func scrapeSeisekiTrRow(s *goquery.Selection) (*model.SeisekiRow, error) {
	var (
		seisekiRow model.SeisekiRow
		err        error
	)

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
		value := strings.ReplaceAll(strings.TrimSpace(elm.AttrOr("title", text)), "　", " ")
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
		case 10:
			// fmt.Printf("i: 10 was ignored\n")
		default:
			err = fmt.Errorf("i: %d row is undefined", i)
			return false
		}
		return true
	})

	return &seisekiRow, err
}
