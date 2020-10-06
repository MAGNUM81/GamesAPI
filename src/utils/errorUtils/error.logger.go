package errorUtils

import (
	"fmt"
	"log"
)

func Fatal(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func Panic(err error) {
	if err != nil {
		str := "PANIC from GamesAPI!"
		panic(fmt.Sprintf("%s\n%s", str, err))
	}
}