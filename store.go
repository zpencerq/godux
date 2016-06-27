package godux

type State interface{}

type Enhancer func(StoreFactory) func(Reducer, State) *Store

type Store struct {
	state State

	reducer       Reducer
	dispatcher    Dispatcher
	isDispatching bool

	listeners     *ListenerSet
	nextListeners *ListenerSet
}

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

func (s *Store) State() State {
	return s.state
}

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

func (s *Store) ReplaceDispatch(d Dispatcher) {
	s.dispatcher = d
}

func (s *Store) Dispatcher() Dispatcher {
	return s.dispatcher
}

func (s *Store) Dispatch(a *Action) *Action {
	return s.dispatcher(a)
}

func (s *Store) ReplaceReducer(r Reducer) {
	s.reducer = r
	s.Dispatch(INITAction())
}
