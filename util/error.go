package util

import "fmt"

func RespStatusError(respCode, expect int) error {
	return fmt.Errorf("Response status was %d(expext %d)", respCode, expect)
}
