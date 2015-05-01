package log

import (
	"fmt"
)

func Info(msg ...interface{}) {
	fmt.Printf(fmt.Sprintf(" %s ", fmt.Sprint(msg...)))
	fmt.Println()

}
