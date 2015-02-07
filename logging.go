package gelf

import (
	"path/filepath"
	"runtime"

	logging "github.com/dkolbly/go-logging"
)

func backend(level logging.Level, calldepth int, rec *logging.Record) (*Common, interface{}) {
	com := &Common{
		ShortMessage: rec.Message(),
		Timestamp:    Time(rec.Time),
		// Note that our logging levels start at CRITICAL=0
		// but the standard syslog has CRITICAL=2, with
		// EMERG and ALERT at levels 0 and 1 respectively
		// and we don't model those
		Level: int(rec.Level) + 2,
	}

	extra := map[string]interface{}{
		"_module": rec.Module,
		"_seq": rec.Id,
	}

	_, file, line, ok := runtime.Caller(calldepth + 1)
	if ok {
		extra["_file"] = filepath.Base(file)
		extra["_line"] = line
	}
	if rec.Annotations != nil {
		for _, annot := range rec.Annotations {
			extra[annot.Key] = annot.Value
		}
	}
	return com, extra
}
