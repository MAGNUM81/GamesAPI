package errorUtils

import "errors"

var (
	ErrNotString   = errors.New("expected value is not a string")
	ErrNoRole      = errors.New("you have no role assigned to you")
	ErrRoleUnknown = errors.New("you have an unknown role assigned to you")
	ErrForbidden   = errors.New("you are not allowed to access specified resource")
)
