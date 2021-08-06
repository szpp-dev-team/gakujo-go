package util

import (
	"fmt"
	"testing"
)

func TestParse2400(t *testing.T) {
	type Pair struct {
		layout string
		value  string
	}
	testcases := []Pair{
		{"2006/01/02 15:04", "2021/07/31 24:00"},
		{"15", "24"},
		{"", ""},
		{"1", "2"},
		{"2006/01/02 15:04:05", "2021/07/31 24:00:00"},
	}
	for _, testcase := range testcases {
		parsedTime, err := Parse2400(testcase.layout, testcase.value)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(parsedTime)
	}
}
