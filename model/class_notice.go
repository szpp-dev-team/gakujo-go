package model

import (
	"time"
)

type ClassNoticeRow struct {
	CourseName      string    //授業科目
	TeachersName string    //担当教員名
	Title        string    //タイトル
	Type         string    //連絡種別
	TargetDate   time.Time //対象日
	Date         time.Time //連絡日時
	Index        int
}
