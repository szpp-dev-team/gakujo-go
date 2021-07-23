package scrape

import "fmt"

type ErrorNotFound struct {
	Name string
}

func (err *ErrorNotFound) Error() string {
	return fmt.Sprintf("%v was not found", err.Name)
}
