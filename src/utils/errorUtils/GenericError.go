package errorUtils

import "fmt"

func NewError(text string) error {
	return fmt.Errorf("%s", text)
}