package godux_test

import (
	. "github.com/zpencerq/godux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BindActionCreators", func() {
	var store *Store
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

	BeforeEach(func() {
		store = NewStore(&StoreInput{Reducer: reducers["todos"]})
	})

	It("Wraps the action creators with the dispatch function", func() {
		actionCreators := map[string]ActionCreator{
			"AddTodo": AddTodo,
		}

		boundActionCreators := BindActionCreators(actionCreators, store.Dispatch)
		Expect(func() {
			_, _ = boundActionCreators["AddTodo"]
		}).To(Not(Panic()))

		action := boundActionCreators["AddTodo"]("Hello")
		Expect(action).To(Equal(actionCreators["AddTodo"]("Hello")))
		Expect(store.State()).To(Equal([]Todo{
			{Id: 1, Text: "Hello"},
		}))
	})

	It("Supports wrapping a single function only", func() {
		actionCreator := AddTodo
		boundActionCreator := BindActionCreator(actionCreator, store.Dispatch)

		action := boundActionCreator("Hello")
		Expect(action).To(Equal(actionCreator("Hello")))
		Expect(store.State()).To(Equal([]Todo{
			{Id: 1, Text: "Hello"},
		}))
	})
})
