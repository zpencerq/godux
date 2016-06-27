package middleware_test

import (
	"reflect"

	. "github.com/zpencerq/godux"
	. "github.com/zpencerq/godux/middleware"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Thunk", func() {
	dispatch := func(action *Action) *Action { return action }
	getState := func() State { return "hi" }
	mc := &MiddlewareContext{getState, dispatch}
	nextHandler := Thunk(mc)

	It("Must return a function to handle next", func() {
		t := reflect.TypeOf(nextHandler(nil)).String()
		Expect(t).To(Equal("godux.Dispatcher"))
	})

	Describe("Handle action", func() {
		It("Must run the given action function with context", func(done Done) {
			actionHandler := nextHandler(nil)

			action := &Action{Type: "thunker", Value: func(mc *MiddlewareContext) *Action {
				close(done)
				return &Action{Type: "ret"}
			}}

			actionHandler(action)
		})

		It("Must pass action to next if not a function", func(done Done) {
			action := &Action{Type: "hello", Value: "world"}

			actionHandler := nextHandler(func(a *Action) *Action {
				Expect(a).To(Equal(action))
				close(done)
				return a
			})

			actionHandler(action)
		})

		It("Must return the return value of next if not a function", func() {
			expected := &Action{Type: "godux"}

			actionHandler := nextHandler(func(a *Action) *Action {
				return expected
			})

			outcome := actionHandler(nil)
			Expect(outcome).To(Equal(expected))
		})

		It("Must return value as expected if a function", func() {
			expected := &Action{Type: "godux"}

			actionHandler := nextHandler(nil)

			outcome := actionHandler(&Action{Value: func(_ *MiddlewareContext) *Action {
				return expected
			}})
			Expect(outcome).To(Equal(expected))
		})

		It("Must be invoked synchronously if a function", func() {
			actionHandler := nextHandler(nil)
			mutated := 0

			actionHandler(&Action{Value: func(_ *MiddlewareContext) *Action {
				mutated++
				return &Action{}
			}})

			Expect(mutated).To(Equal(1))
		})
	})

	Describe("Extra argument", func() {
		It("Must pass the third argument", func(done Done) {
			extraArg := "lol"

			CreateThunk(extraArg)(mc)(nil)(
				&Action{Value: func(ctx *MiddlewareContext, arg interface{}) *Action {
					Expect(mc).To(Equal(ctx))
					Expect(arg).To(Equal(extraArg))
					close(done)
					return &Action{}
				}},
			)
		})
	})
})
