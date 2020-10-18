package errorUtils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewError(text string) error {
	return fmt.Errorf("%s", text)
}

func IsEntityError(c *gin.Context, e EntityError) bool {
	if e != nil {
		c.JSON(e.Status(), e)
	}
	return e != nil
}