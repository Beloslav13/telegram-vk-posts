package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	outFile, _ = os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	LogFile    = log.New(outFile, "", 0)
)

func ForError(err error) bool {
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		dt := time.Now()
		LogFile.Println(fmt.Sprintf("%s: %s", dt.Format("2006-01-02 15:04:05"), err))
		return true
	}
	return false
}
