package log

import (
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/shikaan/kpcli/pkg/errors"
	"github.com/shikaan/kpcli/pkg/info"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FILE = "kpcli.log"
const CONFIG_FOLDER = ".config"

func init() {
	var logPath string
	home, err := os.UserHomeDir()
	if runtime.GOOS == "linux" {
		logPath = strings.Join([]string{home, CONFIG_FOLDER, info.NAME}, string(os.PathSeparator))
	} else {
		logPath = strings.Join([]string{home, info.NAME}, string(os.PathSeparator))
	}

	if err != nil {
		errors.MakeError(err.Error(), "log")
		return
	}

	err = os.MkdirAll(logPath, 0755)

	if err != nil {
		errors.MakeError(err.Error(), "log")
		return
	}

	log.Default().SetOutput(&lumberjack.Logger{
		Filename:   strings.Join([]string{logPath, LOG_FILE}, string(os.PathSeparator)),
		MaxSize:    1,
		MaxBackups: 3,
		MaxAge:     28,
	})
}

func logf(template string, values ...any) {
	log.Default().Printf(template, values...)
}

func Info(msg string) {
	logf("[info] %s", msg)
}

func Error(msg string, err error) {
	logf("[error] %s", msg)
	logf("[debug] %v", err)
}
