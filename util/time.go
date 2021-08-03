package util

import (
	"errors"
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
