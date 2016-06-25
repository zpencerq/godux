package godux_test

import (
	. "github.com/zpencerq/godux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware", func() {
	reducers := map[string]func(interface{}, *Action) interface{}{
		"todos": func(state interface{}, a *Action) interface{} {
			if state == nil {
				state = []Todo{}
			}

			switch a.Type {
			case ADD_TODO:
				slice := state.([]Todo)
				newSlice := make([]Todo, len(state.([]Todo)))
				copy(newSlice, slice)
				newSlice = append(newSlice, Todo{
					Id:   MaxId(slice) + 1,
					Text: a.Value.(string),
				})
				return newSlice
			default:
				return state
			}
		},
	}

	Describe("Apply", func() {
		It("Wraps dispatch method with middleware once", func() {
			fn := func(args ...interface{}) {}
			spy := MakeSpy(fn, &fn)

			test := func(mc *MiddlewareContext) func(Dispatcher) Dispatcher {
				fn(mc)
				return func(dispatcher Dispatcher) Dispatcher {
					return func(action *Action) *Action {
						return dispatcher(action)
					}
				}
			}

			store := Apply(test, Thunk)(NewStore)(&StoreInput{Reducer: reducers["todos"]})

			store.Dispatch(AddTodo("Use Redux"))
			store.Dispatch(AddTodo("Flux FTW!"))

			Expect(spy.Calls).To(HaveLen(1))

			Expect(store.State()).To(Equal([]Todo{
				{Id: 1, Text: "Use Redux"},
				{Id: 2, Text: "Flux FTW!"},
			}))

		})

		It("Passed recursive dispatches through the middleware chain", func(done Done) {
			var nextSpy *Spy

			test := CreateSimpleMiddleware(func(mc *MiddlewareContext, next Dispatcher, action *Action) *Action {
				if nextSpy == nil {
					nextSpy = MakeSpy(next, &next)
				}

				return next(action)
			})

			thunk := CreateSimpleMiddleware(func(mc *MiddlewareContext, next Dispatcher, action *Action) *Action {
				switch a := action.Value.(type) {
				case func(*MiddlewareContext) *Action:
					return a(mc)
				default:
					return next(action)
				}
			})

			store := Apply(test, thunk)(NewStore)(&StoreInput{Reducer: reducers["todos"]})

			go func() {
				if promise, ok := store.Dispatch(AddTodoAsync("Use Redux")).Value.(FutureAction); ok {
					<-promise
					Expect(nextSpy.Calls).To(HaveLen(2))
				}
				close(done)
			}()
		})

		It("Works with thunk middleware", func(done Done) {
			store := Apply(Thunk)(NewStore)(&StoreInput{Reducer: reducers["todos"]})

			store.Dispatch(AddTodoIfEmpty("Hello"))
			Expect(store.State()).To(Equal([]Todo{
				{Id: 1, Text: "Hello"},
			}))

			store.Dispatch(AddTodoIfEmpty("Hello"))
			Expect(store.State()).To(Equal([]Todo{
				{Id: 1, Text: "Hello"},
			}))

			store.Dispatch(AddTodo("World"))
			Expect(store.State()).To(Equal([]Todo{
				{Id: 1, Text: "Hello"},
				{Id: 2, Text: "World"},
			}))

			go func() {
				if promise, ok := store.Dispatch(AddTodoAsync("Maybe")).Value.(FutureAction); ok {
					<-promise
					Expect(store.State()).To(Equal([]Todo{
						{Id: 1, Text: "Hello"},
						{Id: 2, Text: "World"},
						{Id: 3, Text: "Maybe"},
					}))
				}
				close(done)
			}()
		})
	})
})
