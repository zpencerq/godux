package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zpencerq/godux"
	"github.com/zpencerq/godux/middleware"
)

func Increment(i interface{}) *godux.Action {
	return &godux.Action{
		"increment",
		i,
	}
}

func IncrementAsync(i interface{}) *godux.Action {
	return &godux.Action{
		"thunk",
		func(mc *godux.MiddlewareContext) *godux.Action {
			c := make(chan *godux.Action, 1)
			var r <-chan *godux.Action = c
			go func() {
				time.Sleep(1 * time.Second)
				a := mc.Dispatch(Increment(i))
				c <- a
			}()
			return mc.Dispatch(&godux.Action{"promise", r})
		},
	}
}

func Decrement(i interface{}) *godux.Action {
	return &godux.Action{
		"increment",
		i,
	}
}

func Accum(state interface{}, a *godux.Action) interface{} {
	if state == nil {
		state = 0
	}

	val := 1
	if value, ok := a.Value.(int); ok {
		val = value
	}

	switch {
	case a.Type == "increment":
		return state.(int) + val
	case a.Type == "decrement":
		return state.(int) - val
	}
	return state
}

func Type(state interface{}, a *godux.Action) interface{} {
	return a.Type
}

func Identity(state interface{}, a *godux.Action) interface{} {
	return state
}

func main() {

	master := godux.Combine(map[string]godux.Reducer{
		"identity": Identity,
		"k":        Type,
		"accum":    Accum,
	})

	store := godux.Apply(
		middleware.Thunk,
		middleware.LoggingMiddleware,
	)(
		godux.NewStore,
	)(
		&godux.StoreInput{Reducer: master},
	)

	actions := map[string]godux.ActionCreator{
		"Increment":      Increment,
		"IncrementAsync": IncrementAsync,
		"Decrement":      Decrement,
	}
	boundActions := godux.BindActionCreators(actions, store.Dispatch)

	boundActions["Increment"](5)

	promise := NewPromise(func() (interface{}, error) {
		return <-store.Dispatch(IncrementAsync(3)).Value.(<-chan *godux.Action), nil
	})

	time.Sleep(100 * time.Microsecond)

	boundActions["Increment"](1)
	boundActions["Decrement"](5)
	boundActions["Increment"](1)

	fmt.Println("\n============SIMPLE STORE============")

	simpleStore := godux.Apply(
		middleware.CreateLoggingMiddleware(log.New(os.Stdout, "simple: ", log.Flags())),
	)(
		godux.NewStore,
	)(
		&godux.StoreInput{Reducer: Accum},
	)

	simpleBoundActions := godux.BindActionCreators(actions, simpleStore.Dispatch)
	simpleBoundActions["Increment"](1)
	simpleBoundActions["Increment"](3)

	promise.Then(func(v interface{}, _ error) {})
}
