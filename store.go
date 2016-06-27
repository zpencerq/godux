package godux

type State interface{}

type Enhancer func(StoreFactory) func(Reducer, State) *Store

// The only way to change the data in the store is to call Dispatch() on it.
// There should only be a single store in your app. To specify how different
// parts of the state tree respond to actions, you may combine several reducers
// into a single reducer function by using combineReducers.
type Store struct {
	state State

	reducer       Reducer
	dispatcher    Dispatcher
	isDispatching bool

	listeners     *ListenerSet
	nextListeners *ListenerSet
}

// Reducer: A function that returns the next state tree, given
// the current state tree and the action to handle.
//
// State: The initial state. You may optionally specify it
// to hydrate the state from the server in universal apps, or to restore a
// previously serialized user session.
// If you use Combine() to produce the root reducer function, this must be
// an object with the same shape as Combine keys.
//
// Enhancer: The store enhancer. You may optionally specify it
// to enhance the store with third-party capabilities such as middleware,
// time travel, persistence, etc. The only store enhancer that ships with Godux
// is Apply().
type StoreInput struct {
	Reducer  Reducer
	State    State
	Enhancer Enhancer
}

type StoreFactory func(*StoreInput) *Store

func makeDispatcher(s *Store) Dispatcher {
	return func(a *Action) *Action {

		s.listeners = s.nextListeners.Clone()

		if a.Type == "" {
			panic("Action doesn't have a type")
		}

		if s.isDispatching {
			panic("Reducers may not dispatch actions")
		}

		s.isDispatching = true
		defer func() {
			s.isDispatching = false
			s.listeners.Signal()
		}()

		s.state = s.reducer(s.state, a)

		return a
	}
}

// Creates a Godux store that holds the state tree.
//
// When a store is created, an "INIT" action is dispatched so that every
// reducer returns their initial state. This effectively populates
// the initial state tree.
func NewStore(si *StoreInput) *Store {
	if si.Enhancer != nil {
		return si.Enhancer(NewStore)(si.Reducer, si.State)
	}

	store := &Store{
		state:   si.State,
		reducer: si.Reducer,

		listeners:     NewListenerSet(),
		nextListeners: NewListenerSet(),
	}

	store.dispatcher = makeDispatcher(store)

	store.Dispatch(INITAction())

	return store
}

// Reads the state tree managed by the store.
func (s *Store) State() State {
	return s.state
}

// Adds a change listener. It will be called any time an action is dispatched,
// and some part of the state tree may potentially have changed. You may then
// call State() to read the current state tree inside the callback.
//
// You may call Dispatch() from a change listener, with the following
// caveats:
//
// 1. The subscriptions are snapshotted just before every Dispatch() call.
// If you subscribe or unsubscribe while the listeners are being invoked, this
// will not have any effect on the Dispatch() that is currently in progress.
// However, the next Dispatch() call, whether nested or not, will use a more
// recent snapshot of the subscription list.
//
// 2. The listener should not expect to see all state changes, as the state
// might have been updated multiple times during a nested Dispatch() before
// the listener is called. It is, however, guaranteed that all subscribers
// registered before the Dispatch() started will be called with the latest
// state by the time it exits.
func (s *Store) Subscribe(l Listener) func() {
	subscribed := true

	s.nextListeners.Add(&l)

	return func() {
		if !subscribed {
			return
		}

		subscribed = false

		s.nextListeners.Delete(&l)
	}
}

// Dispatches an action. It is the only way to trigger a state change.
//
// The Reducer function, used to create the store, will be called with the
// current state tree and the given action. Its return value will
// be considered the next state of the tree, and the change listeners
// will be notified.
//
// The base implementation only supports plain object actions. If you want to
// dispatch a Promise, an Observable, a thunk, or something else, you need to
// wrap your store creating function into the corresponding middleware. Even the
// middleware will eventually dispatch plain object actions using this method.
//
// Note that, if you use a custom middleware, it may wrap Dispatch() to
// return something else (for example, a Promise you can await).
func (s *Store) Dispatch(a *Action) *Action {
	return s.dispatcher(a)
}

// Replaces the reducer currently used by the store to calculate the state.
//
// You might need this if your app implements code splitting and you want to
// load some of the reducers dynamically.
func (s *Store) ReplaceReducer(r Reducer) {
	s.reducer = r
	s.Dispatch(INITAction())
}
