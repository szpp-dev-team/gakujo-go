package util

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// a wrapper of time.Parse()
func Parse2400(layout, value string) (time.Time, error) {
	parsedTime, err := time.Parse(layout, value)
	if err != nil {
		if !isHourOutErr(err) {
			return time.Time{}, err
		}
		i := strings.Index(layout, "15")
		if i == -1 {
			return time.Time{}, errors.New("stdHour 15 was not found in layout")
		}
		newValue := value[:i] + "00" + value[i+2:]
		parsedTime, err = time.Parse(layout, newValue)
		if err != nil {
			return time.Time{}, err
		}
		return parsedTime.Add(24 * time.Hour), nil
	}
	return parsedTime, nil
}

func isHourOutErr(err error) bool {
	switch err.(type) {
	case *time.ParseError:
		return strings.Contains(err.Error(), "hour")
	default:
		return false
	}
}

func ToWeekday(s rune) time.Weekday {
	switch s {
	case '月':
		return time.Monday
	case '火':
		return time.Tuesday
	case '水':
		return time.Wednesday
	case '木':
		return time.Thursday
	case '金':
		return time.Friday
	case '土':
		return time.Saturday
	case '日':
		return time.Sunday
	default:
		return time.Sunday
	}
}

func BasicTime(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// return beginDate, endDate, error
func ParsePeriod(periodText string) (time.Time, time.Time, error) {
	var beginText1, beginText2, endText1, endText2 string
	fmt.Sscanf(periodText, "%s %s ～ %s %s", &beginText1, &beginText2, &endText1, &endText2)
	beginDate, err := Parse2400("2006/01/02 15:04", fmt.Sprintf("%s %s", beginText1, beginText2))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endDate, err := Parse2400("2006/01/02 15:04", fmt.Sprintf("%s %s", endText1, endText2))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return beginDate, endDate, nil
}
