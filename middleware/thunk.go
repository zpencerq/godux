package middleware

import (
	. "github.com/zpencerq/godux"
)

var Thunk Middleware

func init() {
	Thunk = CreateThunk(nil)
}

func CreateThunk(extraArg interface{}) Middleware {
	return func(mc *MiddlewareContext) func(Dispatcher) Dispatcher {
		return func(next Dispatcher) Dispatcher {
			if next == nil {
				next = mc.Dispatch
			}

			return func(action *Action) *Action {
				if action == nil {
					return next(action)
				}

				switch a := action.Value.(type) {
				case func(*MiddlewareContext) *Action:
					return a(mc)
				case func(*MiddlewareContext, interface{}) *Action:
					return a(mc, extraArg)
				default:
					return next(action)
				}
			}
		}
	}
}
