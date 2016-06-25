package godux

import "reflect"

func Compose(funcs ...interface{}) func(...interface{}) interface{} {
	l := len(funcs)

	if l == 0 {
		return func(i ...interface{}) interface{} {
			return i[0]
		}
	}

	return func(args ...interface{}) interface{} {
		last := funcs[l-1]
		rest := funcs[:l-1]
		r := len(rest)

		in := []reflect.Value{}
		for _, arg := range args {
			in = append(in, reflect.ValueOf(arg))
		}
		out := reflect.ValueOf(last).Call(in)

		for idx, _ := range rest {
			out = reflect.ValueOf(rest[r-idx-1]).Call(out)
		}

		return out[0].Interface()
	}
}
