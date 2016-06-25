package godux_test

import (
	. "github.com/zpencerq/godux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Combine", func() {
	It("Returns a composite reducer that maps the state keys to given reducers", func() {
		reducer := Combine(map[string]Reducer{
			"Counter": Counter,
			"Stack":   Stack,
		})

		s1 := reducer(map[string]interface{}{}, &Action{Type: "increment"})
		Expect(s1).To(Equal(map[string]interface{}{
			"Counter": 1,
			"Stack":   []interface{}{},
		}))
		s2 := reducer(s1, &Action{Type: "push", Value: "a"})
		Expect(s2).To(Equal(map[string]interface{}{
			"Counter": 1,
			"Stack":   []interface{}{"a"},
		}))
	})

	It("Maintains referential equality if the reducers it is combining do", func() {
		reducer := Combine(map[string]Reducer{
			"child1": func(state interface{}, _ *Action) interface{} {
				if _, ok := state.(struct{}); !ok {
					state = struct{}{}
				}
				return state
			},
			"child2": func(state interface{}, _ *Action) interface{} {
				if _, ok := state.(struct{}); !ok {
					state = struct{}{}
				}
				return state
			},
			"child3": func(state interface{}, _ *Action) interface{} {
				if _, ok := state.(struct{}); !ok {
					state = struct{}{}
				}
				return state
			},
		})

		initialState := reducer(nil, INITAction())
		result := reducer(initialState, &Action{Type: "FOO"})
		Expect(&result).To(Equal(&initialState))
	})

	It("Does not have referential equality if one of the reducers changes something", func() {
		reducer := Combine(map[string]Reducer{
			"child1": func(state interface{}, _ *Action) interface{} {
				if _, ok := state.(struct{}); !ok {
					state = struct{}{}
				}
				return state
			},
			"child2": func(state interface{}, action *Action) interface{} {
				if _, ok := state.(struct{ count int }); !ok {
					state = struct{ count int }{0}
				}

				switch action.Type {
				case "increment":
					return struct{ count int }{state.(struct{ count int }).count + 1}
				default:
					return state
				}
			},
			"child3": func(state interface{}, _ *Action) interface{} {
				if _, ok := state.(struct{}); !ok {
					state = struct{}{}
				}
				return state
			},
		})

		initialState := reducer(nil, INITAction())
		result := reducer(initialState, &Action{Type: "increment"})
		Expect(&result).To(Not(Equal(&initialState)))
	})

	It("Throws an error on first call if a reducer attempts to handle a private action", func() {
		reducer := Combine(map[string]Reducer{
			"counter": func(state interface{}, action *Action) interface{} {
				switch action.Type {
				case "increment":
					return state.(int) + 1
				case "decrement":
					return state.(int) - 1

				// Never do this in your code
				case ActionTypes.INIT:
					return 0
				default:
					return state
				}
			},
		})

		Expect(func() { reducer(nil, nil) }).To(Panic())
	})
})

var _ = Describe("CombineReducers", func() {
	It("Returns a composite reducer that maps the state keys to given reducers", func() {
		reducer := CombineReducers(Counter, Stack)

		s1 := reducer(map[string]interface{}{}, &Action{Type: "increment"})
		Expect(s1).To(Equal(map[string]interface{}{
			"Counter": 1,
			"Stack":   []interface{}{},
		}))
		s2 := reducer(s1, &Action{Type: "push", Value: "a"})
		Expect(s2).To(Equal(map[string]interface{}{
			"Counter": 1,
			"Stack":   []interface{}{"a"},
		}))
	})
})
