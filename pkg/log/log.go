package log

import (
	"log"
	"os"
)

func init() {
  f, _ := os.OpenFile("kpcli.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) 
  log.Default().SetOutput(f)
}

func Log(msg ...string) {
  log.Default().Println(msg)
}
func Logf(template string, values ...any) {
  log.Default().Printf(template, values...)
}

