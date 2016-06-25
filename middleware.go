package godux

type MiddlewareContext struct {
	State    func() State
	Dispatch Dispatcher
}

type Middleware func(*MiddlewareContext) func(Dispatcher) Dispatcher

func Apply(middlewares ...Middleware) func(StoreFactory) StoreFactory {
	return func(createStore StoreFactory) StoreFactory {
		return func(si *StoreInput) *Store {
			store := createStore(si)
			dispatch := store.dispatcher

			chain := make([]interface{}, 0, 0)

			context := &MiddlewareContext{
				State: store.State,
				Dispatch: func(a *Action) *Action {
					return dispatch(a)
				},
			}

			for _, middleware := range middlewares {
				chain = append(chain, middleware(context))
			}

			dispatch = Compose(chain...)(dispatch).(Dispatcher)

			store.ReplaceDispatch(dispatch)

			return store
		}
	}
}

func CreateSimpleMiddleware(callback func(*MiddlewareContext, Dispatcher, *Action) *Action) Middleware {
	return func(c *MiddlewareContext) func(Dispatcher) Dispatcher {
		return func(d Dispatcher) Dispatcher {
			return func(a *Action) *Action {
				return callback(c, d, a)
			}
		}
	}
}

func NoopMiddleware(mc *MiddlewareContext) func(Dispatcher) Dispatcher {
	return func(next Dispatcher) Dispatcher {
		return func(action *Action) *Action {
			return next(action)
		}
	}
}
