package model

import (
	"time"
)

type ReportRow struct {
	CourseName     string
	CourseDates    []CourseDate
	Title          string
	Status         Status
	BeginDate      time.Time
	EndDate        time.Time
	LastSubmitDate time.Time
	Format         string
	SubjectMetadata
}

func (rr *ReportRow) DetailOption() *ReportDetailOption {
	return &ReportDetailOption{
		ReportID:        rr.ReportID,
		ListSchoolYear:  rr.SchoolYear,
		ListSubjectCode: rr.SubjectCode,
		ListClassCode:   rr.ClassCode,
		SchoolYear:      rr.SchoolYear,
		SemesterCode:    rr.CourseDates[0].SemesterCode,
	}
}

type SubjectMetadata struct {
	ReportID         string
	SubmitStatusCode string
	SchoolYear       int
	SubjectCode      string
	ClassCode        string
}

type ReportDetail struct {
	Title            string
	BeginDate        time.Time
	EndDate          time.Time
	EvaluationMethod string
	Description      string
	TransMatter      string
}

type Status string

const (
	Accepting = Status("受付中")
	Deadlined = Status("締切")
	Undefined = Status("Undefined")
)

func ToStatus(s string) Status {
	switch s {
	case "受付中":
		return Accepting
	case "締切":
		return Deadlined
	default:
		return Undefined
	}
}

type ReportSearchOption struct {
	SchoolYear   int
	SemesterCode SemesterCode
}

type ReportDetailOption struct {
	ReportID        string
	ListSchoolYear  int
	ListSubjectCode string
	ListClassCode   string
	SchoolYear      int
	SemesterCode    SemesterCode
}
