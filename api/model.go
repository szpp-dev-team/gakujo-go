package api

import (
	"fmt"
)

type TaskType int

const (
	TTMiniTest             TaskType = iota // 小テスト
	TTClassSurvey                          // 授業アンケート
	TTReport                               // レポート
	TTLearningOutcomeSheet                 // 学習成果シート
)

func (tt TaskType) String() string {
	switch tt {
	case TTMiniTest:
		return "小テスト"
	case TTClassSurvey:
		return "授業アンケート"
	case TTReport:
		return "レポート"
	case TTLearningOutcomeSheet:
		return "学習成果シート"
	default:
		return "undefined"
	}
}

func ToTasktype(tt string) (TaskType, error) {
	switch tt {
	case "小テスト":
		return TTMiniTest, nil
	case "授業アンケート":
		return TTClassSurvey, nil
	case "レポート":
		return TTReport, nil
	case "学修成果シート":
		return TTLearningOutcomeSheet, nil
	default:
		return 0, fmt.Errorf("%v is undefined", tt)
	}
}

type NoticeType int

const (
	NTMiniTest      NoticeType = iota // 小テスト
	NTClassSurvey                     // 授業アンケート
	NTCampusSurvey                    // 学内アンケート
	NTReport                          // レポート
	NTCampusContact                   // 学内連絡
	NTClassContact                    // 授業連絡
)

func (nt NoticeType) String() string {
	switch nt {
	case NTMiniTest:
		return "小テスト"
	case NTClassSurvey:
		return "授業アンケート"
	case NTCampusSurvey:
		return "学内アンケート"
	case NTReport:
		return "レポート"
	case NTCampusContact:
		return "学内連絡"
	case NTClassContact:
		return "授業連絡"
	default:
		return "undefined"
	}
}

func ToNoticetype(nt string) (NoticeType, error) {
	switch nt {
	case "小テスト":
		return NTMiniTest, nil
	case "授業アンケート":
		return NTClassSurvey, nil
	case "学内アンケート":
		return NTCampusSurvey, nil
	case "レポート":
		return NTReport, nil
	case "学内連絡":
		return NTCampusContact, nil
	case "授業連絡":
		return NTClassContact, nil
	default:
		return 0, fmt.Errorf("%v is undefined", nt)
	}
}

type SubNoticeType int

const (
	SNTRegist         SubNoticeType = iota // 登録
	SNTTeacherContact                      // 教員連絡
	SNTReminder                            // 催促
	SNTNone                                // なし
)

func (snt SubNoticeType) String() string {
	switch snt {
	case SNTRegist:
		return "登録"
	case SNTTeacherContact:
		return "授業連絡"
	case SNTReminder:
		return "催促"
	default:
		return ""
	}
}

func ToSubNoticetype(snt string) (SubNoticeType, error) {
	switch snt {
	case "登録":
		return SNTRegist, nil
	case "教員連絡":
		return SNTTeacherContact, nil
	case "催促":
		return SNTReminder, nil
	default:
		return 0, fmt.Errorf("%v is undefined", snt)
	}
}
