package middleware

import (
	. "github.com/zpencerq/godux"
)

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
