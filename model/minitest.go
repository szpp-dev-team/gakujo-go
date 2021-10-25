package model

import "time"

type MinitestRow struct {
	CourseName   string
	CourseDates  []CourseDate
	Title        string
	Status       Status
	BeginDate    time.Time
	EndDate      time.Time
	SubmitStatus SubmitStatus
	Format       string
	TaskMetadata
}

func (mr *MinitestRow) DetailOption() *MinitestDetailOption {
	return &MinitestDetailOption{
		TestID:          mr.ID,
		ListSchoolYear:  mr.SchoolYear,
		ListSubjectCode: mr.SubjectCode,
		ListClassCode:   mr.ClassCode,
		SchoolYear:      mr.SchoolYear,
		SemesterCode:    mr.CourseDates[0].SemesterCode,
	}
}

type MinitestDetail struct {
	Title            string
	BeginDate        time.Time
	EndDate          time.Time
	Num              int
	EvaluationMethod string
	Description      string
	TransMatter      string
	MinitestHtml     string
}

type MinitestSearchOption struct {
	SchoolYear   int
	SemesterCode SemesterCode
}

type MinitestDetailOption struct {
	TestID          string
	ListSchoolYear  int
	ListSubjectCode string
	ListClassCode   string
	SchoolYear      int
	SemesterCode    SemesterCode
}

type SubmitStatus string

const (
	UnSubmited = SubmitStatus("未提出")
	Submited   = SubmitStatus("提出済")
)

func ToSubmitStatus(s string) SubmitStatus {
	switch s {
	case "提出済":
		return Submited
	case "未提出":
		return UnSubmited
	}
	return Submited // 適当
}
