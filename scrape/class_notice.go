package scrape

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/szpp-dev-team/gakujo-go/model"
	"github.com/szpp-dev-team/gakujo-go/util"
)

func ClassNoticeRows(r io.Reader) ([]model.ClassNoticeRow, error) {
	parseTime := func(s string) (time.Time, error) {
		if s == "" {
			return time.Time{}, nil
		}
		text := util.ReplaceAndTrim(s)
		t, err := util.Parse2400("2006/01/02 15:04", text)
		if err != nil {
			t, err := time.Parse("2006/01/02", text)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		}
		return t, nil
	}

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	classNoticeRows := []model.ClassNoticeRow{}
	selectror := doc.Find("table#tbl_A01_01 > tbody:nth-child(2) > tr")
	selectror.EachWithBreak(func(i int, s *goquery.Selection) bool {
		classNoticeRow := model.ClassNoticeRow{}
		courseName, courseDate, inerr := parseCourseNameFormat(s.Find("td:nth-child(2)").Text())
		if inerr != nil {
			err = inerr
			return false
		}
		classNoticeRow.CourseName = courseName
		classNoticeRow.CourseDates = courseDate
		classNoticeRow.TeacherName = util.ReplaceAndTrim(s.Find("td:nth-child(3)").Text())
		classNoticeRow.Title = util.ReplaceAndTrim(s.Find("td:nth-child(4) > a").Text())
		classNoticeRow.ContactType = util.ReplaceAndTrim(s.Find("td:nth-child(5)").Text())
		classNoticeRow.TargetDate, inerr = parseTime(s.Find("td:nth-child(6)").Text())
		if inerr != nil {
			err = inerr
			return false
		}
		classNoticeRow.ContactDate, inerr = parseTime(s.Find("td:nth-child(7)").Text())
		if inerr != nil {
			fmt.Println(courseName)
			err = inerr
			return false
		}
		classNoticeRow.Index = i
		classNoticeRows = append(classNoticeRows, classNoticeRow)
		return true
	})
	if err != nil {
		return nil, err
	}
	return classNoticeRows, nil
}

func ClassNoticeDetail(r io.Reader) (*model.ClassNoticeDetail, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	classNoticeDetail := model.ClassNoticeDetail{}
	selector := doc.Find(".ttb_entry > tbody:nth-child(1)")
	classNoticeDetail.ContactType = selector.Find("tr:nth-child(1) > td:nth-child(2)").Text()
	classNoticeDetail.Title = selector.Find("tr:nth-child(2) > td:nth-child(2)").Text()
	classNoticeDetail.Description = selector.Find("tr:nth-child(3) > td:nth-child(2) > div").Text()
	classNoticeDetail.File = nil
	classNoticeDetail.IsFilelinkPublic = false
	classNoticeDetail.ReferenceURL = ""
	classNoticeDetail.Importance = selector.Find("tr:nth-child(7) > td:nth-child(2)").Text()
	classNoticeDetail.ContactDate, err = func() (time.Time, error) {
		text := util.ReplaceAndTrim(selector.Find("tr:nth-child(8) > td:nth-child(2)").Text())
		contactDateText := ""
		fmt.Sscanf(text, "即時通知 %s", &contactDateText)
		return time.Parse("2006/01/02", contactDateText)
	}()
	if err != nil {
		return nil, err
	}
	text := util.ReplaceAndTrim(selector.Find("tr:nth-child(9) > td:nth-child(2)").Text())
	classNoticeDetail.RequireResponse = model.ToRequireResponse(text)
	classNoticeDetail.Importance = util.ReplaceAndTrim(selector.Find("tr:nth-child(10) > td:nth-child(2)").Text())

	return &classNoticeDetail, nil
}

func parseCourseNameFormat(s string) (string, []*model.CourseDate, error) {
	s = strings.TrimSpace(s)
	elems := strings.Split(s, "\n")
	if len(elems) != 2 {
		return "", nil, fmt.Errorf("invalid course name format: ===%s===", s)
	}
	for i := range elems {
		elems[i] = strings.TrimSpace(elems[i])
	}

	courseDates := []*model.CourseDate{}
	for _, plainCourseDate := range strings.Split(elems[1], ",") {
		courseDate, err := parseCourseDateFormat(plainCourseDate)
		if err != nil {
			return "", nil, err
		}
		courseDates = append(courseDates, courseDate)
	}

	return elems[0], courseDates, nil
}

func parseCourseDateFormat(s string) (*model.CourseDate, error) {
	var (
		semester    string
		weekday     rune
		jigen1      int
		jigen2      int
		subSemester string
		other       string
	)

	elms := strings.Split(s, "/")
	if len(elms) != 2 {
		return nil, fmt.Errorf("invalid course date format: ===%s===", s)
	}
	fmt.Sscanf(elms[0], "%s", &semester)

	// 前期/水5・6
	if _, err := fmt.Sscanf(elms[1], "%c%d・%d", &weekday, &jigen1, &jigen2); err != nil {
		if _, err := fmt.Sscanf(elms[1], "%c%d・%d(%s)", &weekday, &jigen1, &jigen2, &subSemester); err != nil {
			if _, err := fmt.Sscanf(elms[1], "%s", &other); err != nil {
				return nil, fmt.Errorf("invalid course date format: ===%s===", s)
			}
		}
	}
	return &model.CourseDate{
		SemesterCode:    model.ToSemesterCode(semester),
		Weekday:         util.ToWeekday(weekday),
		Jigen1:          jigen1,
		Jigen2:          jigen2,
		SubSemesterCode: model.ToSubSemesterCode(subSemester),
		Other:           other,
	}, nil
}
