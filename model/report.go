package model

import (
	"fmt"
	"time"
)

type ReportListInfo struct {
	ReportRows []ReportRow // レポート一覧
}

type ReportRow struct {
	Subject            string       // 科目名
	Title              string       // レポート課題の表示名
	State              ReportState  // レポートの状態
	StartTime          time.Time    // 提出期間の開始日時
	DueTime            time.Time    // 提出期間の終了日時
	LastSubmissionTime time.Time    // 最終提出日時; 未提出の場合はIsZero()がtrue
	Format             ReportFormat // レポート形式
}

type ReportState int

const (
	RSUnknown       ReportState = iota // 未知のレポート状態
	RSAfterDeadline                    // 締切後 (提出済みかは関係なし; 結果未公開)
	RSSubmittable                      // 受付中 (提出済みかは関係なし)
	RSResultOpening                    // 結果公開中
)

type ReportFormat int

const (
	RFUnknown ReportFormat = iota // 未知のレポート形式
	RFWeb                         // Web
	RFPaper                       // 紙
)

func (rs ReportState) String() string {
	switch rs {
	case RSAfterDeadline:
		return "締切"
	case RSSubmittable:
		return "受付中"
	case RSResultOpening:
		return "結果公開中"
	default:
		return "undefined"
	}
}

func ToReportState(rs string) (ReportState, error) {
	switch rs {
	case "締切":
		return RSAfterDeadline, nil
	case "受付中":
		return RSSubmittable, nil
	case "結果公開中":
		return RSResultOpening, nil
	default:
		return RSUnknown, fmt.Errorf("%v is undefined report state", rs)
	}
}

func (rf ReportFormat) String() string {
	switch rf {
	case RFWeb:
		return "Web"
	case RFPaper:
		return "紙"
	default:
		return "undefined"
	}
}

func ToReportFormat(rf string) (ReportFormat, error) {
	switch rf {
	case "Web":
		return RFWeb, nil
	case "紙":
		return RFPaper, nil
	default:
		return RFUnknown, fmt.Errorf("%v is undefined report state", rf)
	}
}
