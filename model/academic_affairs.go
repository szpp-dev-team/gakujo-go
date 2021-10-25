package model

import (
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

type DepartmentGpa struct {
	Gpa            float64
	CalcDate       time.Time // 最終 gpa 算出日
	DepartmentRank int       // 学科内順位
	DepartmentNum  int       // 学科内人数
	CourseRank     int       // コース内順位
	CourseNum      int       // コース内人数
}

type ChusenRegistrationRow struct {
	AttrName           string      // choiceName[i]
	Period             string      // 時限
	SubjectName        string      // 科目名
	ClassName          string      // クラス名
	SubjectDistinction string      // 科目区分
	SubjectType        SubjectType // 必修選択区分
	Credit             int         // 単位
	ChoiceRank         int         // 第n志望 0なら志望なし
	Capacity           int         // 受講定員
	RegistrationStatus             // 履修登録状況
}

type RegistrationStatus struct {
	FirstChoiceNum  int // 第1志望人数
	SecondChoiceNum int // 第2志望人数
	ThirdChoiceNum  int // 第3志望人数
}

type SubjectType string

const (
	STNone               = SubjectType("")
	STCompulsory         = SubjectType("必修")   // 必修
	STElectiveCompulsory = SubjectType("選択必修") // 選択必修
	STElective           = SubjectType("選択")   // 選択
)

func ToSubjectType(st string) SubjectType {
	switch st {
	case "必":
		return SubjectType(STCompulsory)
	case "選択":
		return SubjectType(STElective)
	case "選必":
		return SubjectType(STElectiveCompulsory)
	default:
		return SubjectType(STNone)
	}
}

type GradeType string

const (
	GTFailing   = GradeType("不可") // 不可
	GTFair      = GradeType("可")  // 可
	GTGood      = GradeType("良")  // 良
	GTVerygood  = GradeType("優")  // 優
	GTExcellent = GradeType("秀")  // 秀
	GTPassed    = GradeType("合")  // 合
)

func ToGradeType(gt string) GradeType {
	switch gt {
	case "不可":
		return GradeType(GTFailing)
	case "可":
		return GradeType(GTFair)
	case "良":
		return GradeType(GTGood)
	case "優":
		return GradeType(GTVerygood)
	case "秀":
		return GradeType(GTExcellent)
	case "合":
		return GradeType(GTPassed)
	default:
		return ""
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
