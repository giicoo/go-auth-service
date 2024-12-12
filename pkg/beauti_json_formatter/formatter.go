package beauti_json_formatter

import (
	"bytes"
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	IsColor bool
}

func NewFormatter(IsColor bool) *Formatter {
	return &Formatter{
		IsColor: IsColor,
	}
}

type RootField struct {
	Level      logrus.Level  `json:"level"`
	Msg        string        `json:"msg"`
	Time_stamp string        `json:"time_stamp"`
	Fields     logrus.Fields `json:"fields"`
}

type RootFieldVoid struct {
	Level      logrus.Level `json:"level"`
	Msg        string       `json:"msg"`
	Time_stamp string       `json:"time_stamp"`
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := new(bytes.Buffer)
	var root interface{}

	if len(entry.Data) != 0 {
		root = &RootField{
			Level:      entry.Level,
			Msg:        entry.Message,
			Time_stamp: entry.Time.Format(time.DateTime),
			Fields:     entry.Data,
		}
	} else {
		root = &RootFieldVoid{
			Level:      entry.Level,
			Msg:        entry.Message,
			Time_stamp: entry.Time.Format(time.DateTime),
		}
	}

	data := marshal(root)

	if f.IsColor {
		levelColor := getColorByLevel(entry.Level)
		data = colorize(levelColor, data)
	}

	buffer.WriteString(data)
	buffer.WriteString("\n")
	return buffer.Bytes(), nil
}
