package middleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	. "github.com/zpencerq/godux"
)

var Logging Middleware

func init() {
	Logging = CreateLoggingMiddleware(nil)
}

func CreateLoggingMiddleware(logger *log.Logger) Middleware {
	if logger == nil {
		logger = log.New(os.Stdout, log.Prefix(), log.Flags())
	}

	var long, short bool

	if logger.Flags()&log.Llongfile != 0 {
		long = true
		logger.SetFlags(logger.Flags() - log.Llongfile)
	}
	if logger.Flags()&log.Lshortfile != 0 {
		short = true
		logger.SetFlags(logger.Flags() - log.Lshortfile)
	}

	return CreateSimpleMiddleware(
		func(c *MiddlewareContext, next Dispatcher, action *Action) *Action {
			prefix := ""
			if long || short {
				if _, file, line, ok := runtime.Caller(2); ok {
					if short && !long {
						file = filepath.Base(file)
					}

					prefix = fmt.Sprintf("%s:%d: ", file, line)
				}
			}

			start := time.Now()
			prevState := c.State()

			defer func() {
				took := time.Since(start)
				nextState := c.State()

				logger.Printf("%sAction%v, %v -> %v, %s",
					prefix, *action, prevState, nextState, took)
			}()

			return next(action)
		},
	)
}
