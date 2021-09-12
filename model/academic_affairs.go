package model

import (
	"fmt"
	"net/url"
	"strconv"
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
	Gpa            float64
	TermGpas       []TermGpa
	CalcDate       time.Time // 最終 gpa 算出日
	DepartmentRank int       // 学科内順位
	DepartmentNum  int       // 学科内人数
	CourseRank     int       // コース内順位
	CourseNum      int       // コース内人数
}

type ChusenRegistrationRow struct {
	AttrName           string             // choiceName[i]
	Period             string             // 時限
	SubjectName        string             // 科目名
	ClassName          string             // クラス名
	SubjectDistinction string             // 科目区分
	SubjectType        SubjectType        // 必修選択区分
	Credit             int                // 単位
	ChoiceRank         int                // 第n志望 0なら志望なし
	Capacity           int                // 受講定員
	RegistrationStatus RegistrationStatus // 履修登録状況
}

type RegistrationStatus struct {
	FirstChoiceNum  int // 第1志望人数
	SecondChoiceNum int // 第2志望人数
	ThirdChoiceNum  int // 第3志望人数
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

type PostKamokuFormData struct {
	Faculty       string // 学部番号。学籍番号の最初の2桁？
	Department    string // 学科番号？
	Course        string // コース番号
	Grade         string // 学年
	KamokuKbnCode string // 科目区分コード?
	Req           string // ""
	KamokuCode    string // 科目コード。シラバスでわかると思う
	ClassCode     string // ?
	Unit          string // 単位
	Radio         string // 0-indexed の上から何番目のラジオボタンか
	SelectKamoku  string // radio と同様？
	Youbi         int
	Jigen         int
}

// とりあえず B2 の情報科学科の生徒専用
func NewPostKamokuFormData(kamokuCode, classCode string, unit, radio, youbi, jigen int) *PostKamokuFormData {
	return &PostKamokuFormData{
		Faculty:       "70",
		Department:    "705",
		Course:        "999",
		Grade:         "2",
		KamokuKbnCode: "",
		Req:           "",
		KamokuCode:    kamokuCode,
		ClassCode:     classCode,
		Unit:          strconv.Itoa(unit),
		Radio:         strconv.Itoa(radio),
		SelectKamoku:  strconv.Itoa(radio),
		Youbi:         youbi,
		Jigen:         jigen,
	}
}

func (formData *PostKamokuFormData) FormData() url.Values {
	data := url.Values{}
	data.Set("faculty", formData.Faculty)
	data.Set("department", formData.Department)
	data.Set("course", formData.Course)
	data.Set("grade", formData.Grade)
	data.Set("kamokuKbnCode", formData.KamokuKbnCode)
	data.Set("req", formData.Req)
	data.Set("kamokuCode", formData.KamokuCode)
	data.Set("classCode", formData.ClassCode)
	data.Set("unit", formData.Unit)
	data.Set("radio", formData.Radio)
	data.Set("selectKamoku", formData.SelectKamoku)
	return data

}
