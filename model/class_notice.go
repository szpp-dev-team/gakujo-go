package model

import (
	"net/url"
	"strconv"
	"time"

	"github.com/szpp-dev-team/gakujo-api/util"
)

type ClassNoticeRow struct {
	CourseName  string
	CourseDate  CourseDate
	TeacherName string
	Title       string
	ContactType string
	TargetDate  time.Time
	ContactDate time.Time
	Index       int
}

type CourseDate struct {
	SemesterCode    SemesterCode
	Weekday         time.Weekday
	ComaCsv         string // 1・2コマの授業 = 1,2
	SubSemesterCode SubSemesterCode
	Other           string
}

type ClassNoticeDetail struct {
	ContactType      string          // 連絡種別
	Title            string          // タイトル
	Description      string          // 内容
	File             []string        // 添付ファイル
	IsFilelinkPublic bool            // ファイルリンク公開フラグ
	ReferenceURL     string          // 参照URL
	Importance       string          // 重要度
	ContactDate      time.Time       //連絡日時
	RequireResponse  RequireResponse // WEB返信要求
	Index            int
}

type SemesterCode string

const (
	SemesterCodeNone = SemesterCode("")
	EarlyPeriod      = SemesterCode("前期")
	LaterPeriod      = SemesterCode("後期")
)

func ToSemesterCode(s string) SemesterCode {
	switch s {
	case "前期":
		return EarlyPeriod
	case "後期":
		return LaterPeriod
	default:
		return SemesterCodeNone
	}
}

func (sc SemesterCode) Int() int {
	switch sc {
	case EarlyPeriod:
		return 1
	case LaterPeriod:
		return 2
	case SemesterCodeNone:
		return 0
	default:
		return 0
	}
}

type SubSemesterCode string

const (
	None             = SubSemesterCode("")
	EarlyEarlyPeriod = SubSemesterCode("前期前半")
	EarlyLaterPeriod = SubSemesterCode("前期後半")
	LaterEarlyPeriod = SubSemesterCode("後期前半")
	LaterLaterPeriod = SubSemesterCode("後期後半")
)

func ToSubSemesterCode(s string) SubSemesterCode {
	switch s {
	case "前期前半":
		return EarlyEarlyPeriod
	case "前期後半":
		return EarlyLaterPeriod
	case "後期前半":
		return LaterEarlyPeriod
	case "後期後半":
		return LaterLaterPeriod
	default:
		return ""
	}
}

type ContactKindCode int

const (
	ContactKindCodeNone ContactKindCode = iota
	Canceled
	Supplementary
	Examination
	LectureRoomChange
	TeacherContact
)

func (c *ContactKindCode) String() string {
	switch *c {
	case Canceled, Supplementary, Examination, LectureRoomChange, TeacherContact:
		return strconv.Itoa(int(*c))
	default:
		return ""
	}
}

type RequireResponse bool

const (
	Require    RequireResponse = true
	NotRequire                 = false
)

func ToRequireResponse(s string) RequireResponse {
	if s == "返信を求めない" {
		return NotRequire
	} else {
		return Require
	}
}

func (r RequireResponse) Int() int {
	if r {
		return 1
	}
	return 2
}

type ClassNoticeSearchOption struct {
	TeacherCode                       string
	SchoolYear                        int          // 開講年度
	SemesterCode                      SemesterCode // 開講学期
	SubjectDispCode                   string
	SearchKeyWord                     string          // 検索キーワード
	CheckSearchKeywordTeacherUserName bool            // 担当教員指名
	CheckSearchKeywordSubjectName     bool            // 科目名
	CheckSearchKeywordTitle           bool            // タイトル
	CheckSearchKeywordContent         bool            // 内容
	ContactKindCode                   ContactKindCode // 連絡種別
	TargetDateStart                   time.Time       // 対象日開始
	TargetDateEnd                     time.Time       // 対象日終了
	ReportDateStart                   time.Time       // 連絡日開始
	ReportDateEnd                     time.Time       // 連絡日終了
	RequireResponse                   RequireResponse // WEB返信要求
	OnlyUnRead                        bool            // 未読のみ
	OnlyTodo                          bool            // todoのみ
	OnlyAttachFile                    bool            // 添付ファイルのみ
	StudentCode                       int             // 受信者番号
	StudentName                       string          // 受信者名
}

func BasicClassNoticeSearchOpt(
	year int,
	semesterCode SemesterCode,
	reportDateStart time.Time,
) *ClassNoticeSearchOption {
	return &ClassNoticeSearchOption{
		SchoolYear:                        year,
		SemesterCode:                      semesterCode,
		ReportDateStart:                   reportDateStart,
		CheckSearchKeywordTeacherUserName: true,
		CheckSearchKeywordSubjectName:     true,
		CheckSearchKeywordTitle:           true,
	}
}

func AllClassNoticeSearchOpt(
	year int,
) *ClassNoticeSearchOption {
	return &ClassNoticeSearchOption{
		SchoolYear:                        year,
		ReportDateStart:                   util.BasicTime(2011, 1, 1),
		CheckSearchKeywordTeacherUserName: true,
		CheckSearchKeywordSubjectName:     true,
		CheckSearchKeywordTitle:           true,
	}
}

var whiteList = map[string]struct{}{
	"teacherCode":                       {},
	"schoolYear":                        {},
	"semesterCode":                      {},
	"subjectDispCode":                   {},
	"searchKeyWord":                     {},
	"checkSearchKeywordTeacherUserName": {},
	"checkSearchKeywordSubjectName":     {},
	"checkSearchKeywordTitle":           {},
	"contactKindCode":                   {},
	"targetDateStart":                   {},
	"targetDateEnd":                     {},
	"reportDateStart":                   {},
	"reportDateEnd":                     {},
	"requireResponse":                   {},
	"studentCode":                       {},
	"studentName":                       {},
}

func (o ClassNoticeSearchOption) Formdata() *url.Values {
	on := func(b bool) string {
		if b {
			return "on"
		}
		return ""
	}
	timeText := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.Format("2006/01/02")
	}
	zeroToNone := func(x int) string {
		if x == 0 {
			return ""
		}
		return strconv.Itoa(x)
	}

	data := url.Values{}
	data.Set("teacherCode", o.TeacherCode)
	data.Set("schoolYear", strconv.Itoa(o.SchoolYear))
	data.Set("semesterCode", zeroToNone(o.SemesterCode.Int()))
	data.Set("subjectDispCode", o.SubjectDispCode)
	data.Set("searchKeyWord", o.SearchKeyWord)
	data.Set("checkSearchKeywordTeacherUserName", on(o.CheckSearchKeywordTeacherUserName))
	data.Set("checkSearchKeywordSubjectName", on(o.CheckSearchKeywordSubjectName))
	data.Set("checkSearchKeywordTitle", on(o.CheckSearchKeywordTitle))
	data.Set("checkSearchKeywordContent", on(o.CheckSearchKeywordContent))
	data.Set("contactKindCode", o.ContactKindCode.String())
	data.Set("targetDateStart", timeText(o.TargetDateStart))
	data.Set("targetDateEnd", timeText(o.TargetDateEnd))
	data.Set("reportDateStart", timeText(o.ReportDateStart))
	data.Set("reportDateEnd", timeText(o.ReportDateEnd))
	data.Set("requireResponse", strconv.Itoa(o.RequireResponse.Int()))
	data.Set("onlyUnRead", on(o.OnlyUnRead))
	data.Set("onlyTodo", on(o.OnlyTodo))
	data.Set("onlyAttachFile", on(o.OnlyAttachFile))
	data.Set("studentCode", zeroToNone(o.StudentCode))
	data.Set("studentName", o.StudentName)

	uniqueData := url.Values{}
	for k, v := range data {
		if _, ok := whiteList[k]; ok {
			uniqueData.Set(k, v[0])
		}
		if v[0] != "" {
			uniqueData.Set(k, v[0])
		}
	}

	return &uniqueData
}
