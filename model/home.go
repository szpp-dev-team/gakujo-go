package model

import (
	"time"
)

type HomeInfo struct {
	TaskRows []TaskRow // 未提出課題一覧
}

type TaskRow struct {
	Type     TaskType
	Deadline time.Time
	Name     string
	Index    int
}

type NoticeRow struct {
	Type        NoticeType
	SubType     SubNoticeType
	Important   bool
	Date        time.Time
	Title       string
	Affiliation string
	Index       int
}

type TaskType string

const (
	TTMiniTest             = TaskType("小テスト")    // 小テスト
	TTClassSurvey          = TaskType("授業アンケート") // 授業アンケート
	TTReport               = TaskType("レポート")    // レポート
	TTLearningOutcomeSheet = TaskType("学習成果シート") // 学習成果シート
)

func ToTasktype(tt string) TaskType {
	switch tt {
	case "小テスト":
		return TTMiniTest
	case "授業アンケート":
		return TTClassSurvey
	case "レポート":
		return TTReport
	case "学修成果シート":
		return TTLearningOutcomeSheet
	default:
		return ""
	}
}

type NoticeType string

const (
	NTMiniTest      = NoticeType("小テスト")    // 小テスト
	NTClassSurvey   = NoticeType("授業アンケート") // 授業アンケート
	NTCampusSurvey  = NoticeType("学内アンケート") // 学内アンケート
	NTReport        = NoticeType("レポート")    // レポート
	NTCampusContact = NoticeType("学内連絡")    // 学内連絡
	NTClassContact  = NoticeType("授業連絡")    // 授業連絡
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

func ToNoticetype(nt string) NoticeType {
	switch nt {
	case "小テスト":
		return NTMiniTest
	case "授業アンケート":
		return NTClassSurvey
	case "学内ｱﾝｹｰﾄ": // これ頭おかしいよ・・・
		return NTCampusSurvey
	case "レポート":
		return NTReport
	case "学内連絡":
		return NTCampusContact
	case "授業連絡":
		return NTClassContact
	default:
		return ""
	}
}

type SubNoticeType string

const (
	SNTRegist            = SubNoticeType("登録")    // 登録
	SNTTeacherContact    = SubNoticeType("教員連絡")  // 教員連絡
	SNTReminder          = SubNoticeType("催促")    // 催促
	SNTComment           = SubNoticeType("コメント")  // コメント
	SNTChangeLectureRoom = SubNoticeType("講義室変更") // 講義室変更
	SNTNone                                       // なし
)

func ToSubNoticetype(snt string) SubNoticeType {
	switch snt {
	case "登録":
		return SNTRegist
	case "教員連絡":
		return SNTTeacherContact
	case "催促":
		return SNTReminder
	case "コメント":
		return SNTComment
	case "講義室変更":
		return SNTChangeLectureRoom
	default:
		return SNTNone
	}
}

type NoticeDetail struct {
	ContactType         string    //連絡種別
	Title               string    //タイトル
	Detail              string    //連絡内容
	File                string    // ファイル
	FilelinkPublication bool      //ファイルリンク公開
	ReferenceURL        string    //参照URL
	Important           bool      //重要度
	Date                time.Time //日時
	WebReturnRequest    bool      //WEB返信要求
}
