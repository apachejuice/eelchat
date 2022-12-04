package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/inconshreveable/log15"
	"github.com/spf13/viper"
)

var levelMap map[log15.Lvl]string = map[log15.Lvl]string{
	log15.LvlCrit:  "critical",
	log15.LvlError: "error",
	log15.LvlWarn:  "warn",
	log15.LvlInfo:  "info",
	log15.LvlDebug: "debug",
}

type Logger interface {
	log15.Logger
	Fatal(msg string, ctx ...any)
}

type loggerWrapper struct {
	log15.Logger
}

func (l *loggerWrapper) Fatal(msg string, ctx ...any) {
	l.Crit(fmt.Sprintf("Fatal error, exiting: %s", msg), ctx...)
	log.Fatal(msg)
}

// Returns a new logger instance for the given component name.
// log15 uses rather unconventional log records, so we have to mutate the output somewhat.
func NewLogger(component string) Logger {
	logger := log15.New("component", component)
	logFile := viper.GetString("logging.file") // no constant due to import cycle
	formatFunc := log15.FormatFunc(func(r *log15.Record) []byte {
		obj := map[string]string{
			"level":     levelMap[r.Lvl],
			"timestamp": r.Time.String(),
			"message":   r.Msg,
		}

		extras := r.Ctx
		if len(extras)%2 != 0 {
			extras = append(extras, nil)
		}

		for i := 0; i < len(extras); i += 2 {
			obj[fmt.Sprintf("%s", extras[i])] = fmt.Sprintf("%s", extras[i+1])
		}

		rec, _ := json.Marshal(obj)
		return append(rec, '\n')
	})

	fhandler := log15.Must.FileHandler(logFile, formatFunc)
	if os.Getenv("EELCHAT_DEBUG") != "" {
		fhandler = log15.CallerFileHandler(fhandler)
	}

	logger.SetHandler(log15.MultiHandler(
		fhandler, log15.LvlFilterHandler(
			log15.LvlError, log15.StreamHandler(os.Stderr, log15.TerminalFormat()),
		)),
	)

	return &loggerWrapper{Logger: logger}
}
