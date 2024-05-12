package log

import (
	"bytes"
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
)

type Formatter string

const pathKey = "PATH"

func (f Formatter) Format(entry *log.Entry) ([]byte, error) {
	if entry.Buffer == nil {
		entry.Buffer = &bytes.Buffer{}
	}
	sb := NewPrinterOn(entry.Buffer)
	dm := entry.Data
	if path, ok := dm[pathKey].(string); ok {
		sb.Printf("%s %-*s %s : %s", entry.Time.Format(string(f)), len("warning"), entry.Level, path, entry.Message)
	} else {
		sb.Printf("%s %-*s %s", entry.Time.Format(string(f)), len("warning"), entry.Level, entry.Message)
	}
	printFields(dm, sb)
	sb.PrintByte('\n')
	return entry.Buffer.Bytes(), nil
}

type PlainFormatter struct{}

func (PlainFormatter) Format(entry *log.Entry) ([]byte, error) {
	if entry.Buffer == nil {
		entry.Buffer = &bytes.Buffer{}
	}
	sb := Printer{entry.Buffer}
	sb.PrintString(entry.Message)
	printFields(entry.Data, sb)
	sb.PrintByte('\n')
	return entry.Buffer.Bytes(), nil
}

func printFields(dm log.Fields, sb Printer) {
	if len(dm) > 0 {
		ks := make([]string, 0, len(dm))
		for k := range dm {
			if k != pathKey {
				ks = append(ks, k)
			}
		}
		sort.Slice(ks, func(i, j int) bool {
			return ks[i] < ks[j]
		})
		sb.PrintString(" :")
		for _, k := range ks {
			sb.PrintByte(' ')
			sb.PrintString(k)
			sb.PrintByte('=')
			v := dm[k]
			switch v := v.(type) {
			case string, fmt.Stringer:
				sb.Printf("%q", v)
			default:
				sb.Printf("%v", v)
			}
		}
	}
}
