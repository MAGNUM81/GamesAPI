package errors

import "fmt"

func New(text string) error {
	return fmt.Errorf("%s", text)
}