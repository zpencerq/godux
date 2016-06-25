package main

type Promise func() (interface{}, error)

func (p Promise) Then(callback func(interface{}, error)) {
	ret, err := p()
	callback(ret, err)
}

func NewPromise(f func() (interface{}, error)) Promise {
	var result interface{}
	var err error

	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()

	return func() (interface{}, error) {
		<-c
		return result, err
	}
}

func SimplePromise(f func() interface{}) func() interface{} {
	return func() interface{} {
		p, _ := NewPromise(func() (interface{}, error) {
			return f(), nil
		})()
		return p
	}
}
