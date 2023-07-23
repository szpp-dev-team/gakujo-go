package model

import "time"

type ClassEnqRow struct {
	CourseName   string
	CourseDates  []*CourseDate
	Title        string
	Status       Status
	BeginDate    time.Time
	EndDate      time.Time
	SubmitStatus SubmitStatus
	TaskMetadata
}

func (cr *ClassEnqRow) DetailOption() *ClassEnqDetailOption {
	return &ClassEnqDetailOption{
		ClassEnqID:      cr.ID,
		ListSchoolYear:  cr.SchoolYear,
		ListSubjectCode: cr.SubjectCode,
		ListClassCode:   cr.ClassCode,
		SchoolYear:      cr.SchoolYear,
		SemesterCode:    cr.CourseDates[0].SemesterCode,
	}
}

type ClassEnqDetail struct {
	Title        string
	BeginDate    time.Time
	EndDate      time.Time
	Num          int
	NameType     string
	Description  string
	TransMatter  string
	ClassEnqHtml string
}

type ClassEnqSearchOption struct {
	SchoolYear   int
	SemesterCode SemesterCode
}

type ClassEnqDetailOption struct {
	ClassEnqID      string
	ListSchoolYear  int
	ListSubjectCode string
	ListClassCode   string
	SchoolYear      int
	SemesterCode    SemesterCode
}
