package scrape

import (
	"testing"
	"time"
)

func TestParseReportSubmissionPeriod(t *testing.T) {
	jst, _ := time.LoadLocation("Asia/Tokyo")

	cases := []struct {
		input string
		start time.Time
		due   time.Time
	}{
		{
			input: "2021/06/29 00:00 ～ 2021/08/03 19:00",
			start: time.Date(2021, 6, 29, 0, 0, 0, 0, jst),
			due:   time.Date(2021, 8, 3, 19, 0, 0, 0, jst),
		},
		{
			input: "2021/12/31 24:00 ～ 2022/01/01 00:00",
			start: time.Date(2022, 1, 1, 0, 0, 0, 0, jst),
			due:   time.Date(2022, 1, 1, 0, 0, 0, 0, jst),
		},
	}

	for _, c := range cases {
		start, due, err := parseReportSubmissionPeriod(c.input)
		if err != nil || !start.Equal(c.start) || !due.Equal(c.due) {
			t.Fatalf("input: %q\n"+
				"got:      (%q %q %q)\n"+
				"expected: (%q %q)",
				c.input, start, due, err, c.start, c.due)
		}
	}
}
