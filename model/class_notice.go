package model

import (
	"time"
)

type ClassNoticeRow struct {
	CourseName   string    `json:"course_name,omitempty"`   //授業科目
	TeachersName string    `json:"teachers_name,omitempty"` //担当教員名
	Title        string    `json:"title,omitempty"`         //タイトル
	Type         string    `json:"type,omitempty"`          //連絡種別
	TargetDate   time.Time `json:"target_date,omitempty"`   //対象日
	Date         time.Time `json:"date,omitempty"`          //連絡日時
	Index        int       `json:"index,omitempty"`
}
