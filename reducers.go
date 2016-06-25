package godux

import (
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Reducer func(interface{}, *Action) interface{}

type undef struct {
	__can_not_use int
}

func assertReducerSanity(reducers map[string]Reducer, final Reducer) error {
	for name, reducer := range reducers {
		undefined := undef{3}

		initState := reducer(undefined, INITAction())

		t := "@@godux/PROBE_UNKNOWN_ACTION" + strings.Join(strings.Split(strconv.FormatInt(rand.Int63(), 36)[:7], ""), ".")
		probeState := reducer(undefined, &Action{Type: t})

		_, initOk := initState.(undef)
		_, probeOk := probeState.(undef)
		if (!initOk && probeOk) || (initOk && !probeOk) {
			return errors.New(
				`Reducer ` + name + ` differed between init action and the
random probe. Don't try to handle ` + ActionTypes.INIT + `
or other actions in the "@@godux/*" namespace.
They are considered private. Instead, you must return the
current state for any unknown actions.`)
		}
	}

	return nil
}

func Combine(reducers map[string]Reducer) Reducer {
	var sanityError error

	finalReducer := func(previousState interface{}, a *Action) interface{} {
		if sanityError != nil {
			panic(sanityError)
		}

		if previousState == nil {
			previousState = make(map[string]interface{})
		}
		hasChanged := false
		nextState := make(map[string]interface{})

		for key, reducer := range reducers {
			previousValue, _ := previousState.(map[string]interface{})[key]
			nextValue := reducer(previousValue, a)
			nextState[key] = nextValue
			hasChanged = hasChanged || (&nextValue != &previousValue)
		}

		if hasChanged {
			return nextState
		} else {
			return previousState
		}
	}

	sanityError = assertReducerSanity(reducers, finalReducer)

	return finalReducer
}

func CombineReducers(reducers ...Reducer) Reducer {
	reducerMap := map[string]Reducer{}
	for _, reducer := range reducers {
		fullName := runtime.FuncForPC(reflect.ValueOf(reducer).Pointer()).Name()
		s := strings.Split(fullName, ".")
		reducerMap[s[len(s)-1]] = reducer
	}

	return func(previousState interface{}, a *Action) interface{} {
		if previousState == nil {
			previousState = make(map[string]interface{})
		}
		hasChanged := false
		nextState := make(map[string]interface{})

		for key, reducer := range reducerMap {
			previousValue, _ := previousState.(map[string]interface{})[key]
			nextValue := reducer(previousValue, a)
			nextState[key] = nextValue
			hasChanged = hasChanged || (&nextValue != &previousValue)
		}

		if hasChanged {
			return nextState
		} else {
			return previousState
		}
	}
}
