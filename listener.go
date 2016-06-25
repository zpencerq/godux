package godux

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
)

type Listener func()

type ListenerSet struct {
	*sync.RWMutex
	*list.List
}

func NewListenerSet() *ListenerSet {
	return &ListenerSet{&sync.RWMutex{}, list.New()}
}

func (ls *ListenerSet) Equal(o *ListenerSet) bool {
	ls.RLock()
	o.RLock()
	defer ls.RUnlock()
	defer o.RUnlock()

	if ls.Len() != o.Len() {
		return false
	}

	e := ls.Front()
	f := o.Front()
	for e != nil && f != nil {
		if e != f {
			return false
		}
		if e != nil {
			e = e.Next()
		}
		if f != nil {
			f = f.Next()
		}
	}

	return true
}

func (ls *ListenerSet) Add(l *Listener) *list.Element {
	ls.Lock()
	defer ls.Unlock()

	return ls.PushBack(l)
}

func (ls *ListenerSet) Delete(listener *Listener) *Listener {
	ls.Lock()
	defer ls.Unlock()

	for e := ls.Front(); e != nil; e = e.Next() {
		if e.Value == listener {
			return ls.Remove(e).(*Listener)
		}
	}
	return nil
}

func (ls *ListenerSet) Clone() *ListenerSet {
	ls.RLock()
	defer ls.RUnlock()

	listeners := NewListenerSet()
	for e := ls.Front(); e != nil; e = e.Next() {
		listeners.PushBack(e.Value)
	}
	return listeners
}

func (ls *ListenerSet) Signal() {
	ls.RLock()
	defer ls.RUnlock()

	for e := ls.Front(); e != nil; e = e.Next() {
		(*e.Value.(*Listener))()
	}
}

func (ls *ListenerSet) String() string {
	ls.RLock()
	defer ls.RUnlock()

	result := make([]string, ls.Len())
	for e := ls.Front(); e != nil; e = e.Next() {
		result = append(result, fmt.Sprintf("%s", e.Value))
	}
	return strings.Join(result, ",")
}
