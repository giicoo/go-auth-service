package beauti_json_formatter

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return lightGray
	case logrus.DebugLevel:
		return lightBlue
	case logrus.InfoLevel:
		return lightGreen
	case logrus.WarnLevel:
		return lightYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return lightRed
	default:
		return blue
	}
}
