package middleware

import (
	"log"
	"os"
	"time"

	. "github.com/zpencerq/godux"
)

func CreateLoggingMiddleware(logger *log.Logger) Middleware {
	if logger == nil {
		logger = log.New(os.Stdout, log.Prefix(), log.Flags())
	}

	return CreateSimpleMiddleware(
		func(c *MiddlewareContext, next Dispatcher, action *Action) *Action {
			start := time.Now()
			prevState := c.State()

			defer func() {
				took := time.Since(start)
				nextState := c.State()

				logger.Printf("Action%v, %v -> %v, %s",
					*action, prevState, nextState, took)
			}()

			return next(action)
		},
	)
}

func LoggingMiddleware(mc *MiddlewareContext) func(Dispatcher) Dispatcher {
	return func(dispatcher Dispatcher) Dispatcher {
		return func(action *Action) *Action {
			start := time.Now()
			prevState := mc.State()

			defer func() {
				took := time.Since(start)
				nextState := mc.State()

				log.Printf("Action%v, %v -> %v, %s",
					*action, prevState, nextState, took)
			}()

			return dispatcher(action)
		}
	}
}
