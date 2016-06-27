package middleware_test

import (
	"reflect"
)

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
