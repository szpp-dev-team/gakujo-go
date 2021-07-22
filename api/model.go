package api

import (
	"fmt"
)

type TaskType int

const (
	MiniTest             TaskType = iota // 小テスト
	ClassSurvey                          // 授業アンケート
	Report                               // レポート
	LearningOutcomeSheet                 // 学習成果シート
)

func (tt TaskType) String() string {
	switch tt {
	case MiniTest:
		return "小テスト"
	case ClassSurvey:
		return "授業アンケート"
	case Report:
		return "レポート"
	case LearningOutcomeSheet:
		return "学習成果シート"
	default:
		return "undefined"
	}
}

func ToTasktype(tt string) (TaskType, error) {
	switch tt {
	case "小テスト":
		return MiniTest, nil
	case "授業アンケート":
		return ClassSurvey, nil
	case "レポート":
		return Report, nil
	case "学修成果シート":
		return LearningOutcomeSheet, nil
	default:
		return 0, fmt.Errorf("%v is undefined", tt)
	}
}
