package godux_test

import (
	"reflect"

	. "github.com/zpencerq/godux"
)

type Todo struct {
	Id   int
	Text string
}

var ADD_TODO string = "ADD_TODO"

func AddTodo(text interface{}) *Action {
	return &Action{Type: ADD_TODO, Value: text}
}

func UnknownAction() *Action {
	return &Action{Type: "UNKNOWN"}
}

func MaxId(todos []Todo) int {
	max := 0

	for _, todo := range todos {
		if todo.Id > max {
			max = todo.Id
		}
	}

	return max
}

type Call struct {
	In  []reflect.Value
	Out []reflect.Value
}

type Spy struct {
	Fn    interface{}
	Calls []*Call
}

func (s *Spy) Wrap(f interface{}) interface{} {
	v := reflect.ValueOf(f)
	return reflect.MakeFunc(v.Type(),
		func(in []reflect.Value) []reflect.Value {
			out := v.Call(in)
			s.Calls = append(s.Calls, &Call{in, out})
			return out
		},
	)
}

func (s *Spy) Call(in []reflect.Value) []reflect.Value {
	return reflect.ValueOf(s.Fn).Call(in)
}

func (s *Spy) WasCalledWith(values ...interface{}) bool {
	for _, call := range s.Calls {
		result := true
		for idx, arg := range call.In {
			if values[idx] != arg.Interface() {
				result = false
			}
		}
		if result {
			return result
		}
	}
	return true
}

func MakeSpy(fn interface{}, fptr interface{}) *Spy {
	spy := &Spy{Calls: make([]*Call, 0)}

	if fptr != nil {
		fptrE := reflect.ValueOf(fptr).Elem()
		fptrE.Set(spy.Wrap(fn).(reflect.Value))
	} else {
		spy.Fn = spy.Wrap(fn)
	}

	return spy
}

func Foo(state interface{}, action *Action) interface{} {
	if _, ok := state.(int); !ok {
		state = 0
	}

	switch {
	case action.Type == "foo":
		return 1
	default:
		return state
	}
}

func Bar(state interface{}, action *Action) interface{} {
	if _, ok := state.(int); !ok {
		state = 0
	}

	switch {
	case action.Type == "bar":
		return 2
	default:
		return state
	}
}

func Counter(state interface{}, action *Action) interface{} {
	if _, ok := state.(int); !ok {
		state = 0
	}

	if action.Type == "increment" {
		return state.(int) + 1
	} else {
		return state
	}
}

func Stack(state interface{}, action *Action) interface{} {
	if _, ok := state.([]interface{}); !ok {
		state = []interface{}{}
	}

	if action.Type == "push" {
		return append(state.([]interface{}), action.Value)
	} else {
		return state
	}
}

type ThunkFunc func(*MiddlewareContext) *Action

type FutureAction <-chan *Action

func AddTodoAsync(text string) *Action {
	return &Action{Type: ADD_TODO, Value: func(mc *MiddlewareContext) *Action {
		return mc.Dispatch(AddTodo(text))
	}}
}

func AddTodoIfEmpty(text string) *Action {
	return &Action{Type: ADD_TODO, Value: func(mc *MiddlewareContext) *Action {
		if len(mc.State().([]Todo)) == 0 {
			return mc.Dispatch(AddTodo(text))
		} else {
			return nil
		}
	}}
}

func Thunk(mc *MiddlewareContext) func(Dispatcher) Dispatcher {
	return func(next Dispatcher) Dispatcher {
		return func(action *Action) *Action {
			switch a := action.Value.(type) {
			case func(*MiddlewareContext) *Action:
				return a(mc)
			default:
				return next(action)
			}
		}
	}
}

var Reducers map[string]Reducer

func init() {
	Reducers = map[string]Reducer{
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
		"todosReverse": func(state interface{}, a *Action) interface{} {
			if state == nil {
				state = []Todo{}
			}

			switch a.Type {
			case ADD_TODO:
				slice := state.([]Todo)
				newSlice := make([]Todo, len(state.([]Todo)))
				copy(newSlice, slice)
				newSlice = append([]Todo{Todo{
					Id:   MaxId(slice) + 1,
					Text: a.Value.(string),
				}}, slice...)
				return newSlice
			default:
				return state
			}
		},
	}
}
