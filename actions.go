package godux

import "reflect"

type Action struct {
	Type  string
	Value interface{}
}

type Dispatcher func(*Action) *Action

type ActionCreator func(interface{}) *Action

var ActionTypes struct{ INIT string }

func INITAction() *Action {
	return &Action{Type: ActionTypes.INIT}
}

func init() {
	ActionTypes = struct{ INIT string }{
		INIT: "@@godux/INIT",
	}
}

func BindActionCreator(creator ActionCreator, dispatch Dispatcher) ActionCreator {
	return func(i interface{}) *Action {
		return dispatch(creator(i))
	}
}

func BindActionCreators(creators map[string]ActionCreator, dispatch Dispatcher) map[string]ActionCreator {
	result := map[string]ActionCreator{}
	for k, creator := range creators {
		result[k] = BindActionCreator(creator, dispatch)
	}
	return result
}

// This replaces the fields in the struct with bound creators
func BindActionCreatorsStruct(creators interface{}, dispatch Dispatcher) interface{} {
	value := reflect.ValueOf(creators)
	if value.Elem().Kind() != reflect.Struct {
		panic("Creators not given in *struct")
	}

	copy := creators

	for idx := 0; idx < value.Elem().NumField(); idx++ {
		field := reflect.ValueOf(creators).Elem().Field(idx)

		fptr := field.Interface().(ActionCreator)
		wrapper := func(i interface{}) *Action {
			return dispatch(fptr(i))
		}

		dstField := reflect.ValueOf(copy).Elem().Field(idx)
		ptr := dstField.Addr().Interface().(*ActionCreator)
		*ptr = wrapper
	}
	return copy
}
