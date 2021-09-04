package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type SeisekiRow struct {
	SubjectName        string      `json:"subject_name,omitempty"`        // 科目名
	TeacherName        string      `json:"teacher_name,omitempty"`        // 教員名
	SubjectDistinction string      `json:"subject_distinction,omitempty"` // 科目区分
	SubjectType        SubjectType `json:"subject_type,omitempty"`        // 必修選択区分
	Credit             int         `json:"credit,omitempty"`              // 単位
	Grade              GradeType   `json:"grade,omitempty"`               // 評価
	Score              float64     `json:"score,omitempty"`               // 得点
	GP                 float64     `json:"gp,omitempty"`                  // 科目GP
	Year               int         `json:"year,omitempty"`                // 取得年度
	Date               time.Time   `json:"date,omitempty"`                // 報告日
	//TestType           string      // 試験種別 よくわからないので保留
}

type TermGpa struct {
	Year      int     `json:"year,omitempty"`
	FirstGPA  float64 `json:"first_gpa,omitempty"`  // 前期 gpa
	SecondGPA float64 `json:"second_gpa,omitempty"` // 後期 gpa
}

type DepartmentGpa struct {
	Gpa            float64   `json:"gpa,omitempty"` // 累積 gpa
	TermGpas       []TermGpa `json:"term_gpas,omitempty"`
	CalcDate       time.Time `json:"calc_date,omitempty"`       // 最終 gpa 算出日
	DepartmentRank int       `json:"department_rank,omitempty"` // 学科内順位
	DepartmentNum  int       `json:"department_num,omitempty"`  // 学科内人数
	CourseRank     int       `json:"course_rank,omitempty"`     // コース内順位
	CourseNum      int       `json:"course_num,omitempty"`      // コース内人数
}

type ChusenRegistrationRow struct {
	AttrName           string             `json:"attr_name,omitempty"`           // choiceName[i]
	Period             string             `json:"period,omitempty"`              // 時限
	SubjectName        string             `json:"subject_name,omitempty"`        // 科目名
	ClassName          string             `json:"class_name,omitempty"`          // クラス名
	SubjectDistinction string             `json:"subject_distinction,omitempty"` // 科目区分
	SubjectType        SubjectType        `json:"subject_type,omitempty"`        // 必修選択区分
	Credit             int                `json:"credit,omitempty"`              // 単位
	ChoiceRank         int                `json:"choice_rank,omitempty"`         // 第n志望 0なら志望なし
	Capacity           int                `json:"capacity,omitempty"`            // 受講定員
	RegistrationStatus RegistrationStatus `json:"registration_status,omitempty"` // 履修登録状況
}

type RegistrationStatus struct {
	FirstChoiceNum  int `json:"first_choice_num,omitempty"`  // 第1志望人数
	SecondChoiceNum int `json:"second_choice_num,omitempty"` // 第2志望人数
	ThirdChoiceNum  int `json:"third_choice_num,omitempty"`  // 第3志望人数
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

func (st SubjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(st.String())
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

func (gt GradeType) MarshalJSON() ([]byte, error) {
	return json.Marshal(gt.String())
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
