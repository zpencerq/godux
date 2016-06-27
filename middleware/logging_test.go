package middleware_test

import (
	"bytes"
	"log"
	"path/filepath"
	"runtime"

	. "github.com/zpencerq/godux"
	. "github.com/zpencerq/godux/middleware"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	dispatch := func(action *Action) *Action { return action }
	getState := func() State { return "hi" }
	mc := &MiddlewareContext{getState, dispatch}

	It("Logs when an action is dispatched", func() {
		var buf bytes.Buffer
		logger := log.New(&buf, "logger: ", 0)

		logMiddleware := CreateLoggingMiddleware(logger)
		action := &Action{Type: "hello", Value: "world"}
		logMiddleware(mc)(nil)(action)

		Expect(buf.String()).To(MatchRegexp(`logger: Action{hello world}, hi -> hi, `))
	})

	Context("Shows proper call location", func() {
		It("For os.Log's log.Lshortfile flag", func() {
			var buf bytes.Buffer
			logger := log.New(&buf, "logger: ", log.Lshortfile)

			logMiddleware := CreateLoggingMiddleware(logger)
			action := &Action{Type: "hello", Value: "world"}
			logMiddleware(mc)(nil)(action)

			_, file, line, _ := runtime.Caller(0)
			Expect(buf.String()).To(
				MatchRegexp(`logger: %s:%d: Action{hello world}, hi -> hi, `,
					filepath.Base(file), line-2),
			)
		})

		It("For os.Log's log.Llongfile flag", func() {
			var buf bytes.Buffer
			logger := log.New(&buf, "logger: ", log.Llongfile)

			logMiddleware := CreateLoggingMiddleware(logger)
			action := &Action{Type: "hello", Value: "world"}
			logMiddleware(mc)(nil)(action)

			_, file, line, _ := runtime.Caller(0)
			Expect(buf.String()).To(
				MatchRegexp(`logger: %s:%d: Action{hello world}, hi -> hi, `,
					file, line-2),
			)
		})
	})
})
