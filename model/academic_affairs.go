package model

import (
	"fmt"
	"time"
)

type SeisekiRow struct {
	SubjectName        string      // 科目名
	TeacherName        string      // 教員名
	SubjectDistinction string      // 科目区分
	SubjectType        SubjectType // 必修選択区分
	Credit             int         // 単位
	Grade              GradeType   // 評価
	Score              float64     // 得点
	GP                 float64     // 科目GP
	Year               int         // 取得年度
	Date               time.Time   // 報告日
	//TestType           string      // 試験種別 よくわからないので保留
}

type TermGpa struct {
	Year      int
	FirstGPA  float64 // 前期 gpa
	SecondGPA float64 // 後期 gpa
}

type DepartmentGpa struct {
	Gpa            float64 // 累積 gpa
	TermGpas       []TermGpa
	CalcDate       time.Time // 最終 gpa 算出日
	DepartmentRank int       // 学科内順位
	DepartmentNum  int       // 学科内人数
	CourseRank     int       // コース内順位
	CourseNum      int       // コース内人数
}

type SubjectType int

const (
	STCompulsory         SubjectType = iota // 必修
	STElectiveCompulsory                    // 選択必修
	STElective                              // 選択
)

func (st SubjectType) String() string {
	switch st {
	case STCompulsory:
		return "必修"
	case STElective:
		return "選択"
	case STElectiveCompulsory:
		return "選択必修"
	default:
		return "undefined"
	}
}

func ToSubjectType(st string) (SubjectType, error) {
	switch st {
	case "必":
		return SubjectType(STCompulsory), nil
	case "選択":
		return SubjectType(STElective), nil
	case "選必":
		return SubjectType(STElectiveCompulsory), nil
	default:
		return 0, fmt.Errorf("%v is undefined", st)
	}
}

type GradeType int

// TODO: ガバガバ英語をなんとかする
const (
	GTFailing   GradeType = iota // 不可
	GTFair                       // 可
	GTGood                       // 良
	GTVerygood                   // 優
	GTExcellent                  // 秀
	GTPassed                     // 合
)

func (gt GradeType) String() string {
	switch gt {
	case GTFailing:
		return "不可"
	case GTFair:
		return "可"
	case GTGood:
		return "良"
	case GTVerygood:
		return "優"
	case GTExcellent:
		return "秀"
	case GTPassed:
		return "合"
	default:
		return "undefined"
	}
}

func ToGradeType(gt string) (GradeType, error) {
	switch gt {
	case "不可":
		return GradeType(GTFailing), nil
	case "可":
		return GradeType(GTFair), nil
	case "良":
		return GradeType(GTGood), nil
	case "優":
		return GradeType(GTVerygood), nil
	case "秀":
		return GradeType(GTExcellent), nil
	case "合":
		return GradeType(GTPassed), nil
	default:
		return 0, fmt.Errorf("%v is undefined", gt)
	}
}
