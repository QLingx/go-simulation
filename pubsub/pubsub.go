package pubsub

import (
	"fmt"
	"runtime"
	"reflect"
	"sync"
	"errors"
)

type PubSubError struct{
	f interface{}
	e interface{}
}

func (pse *PubSubError)String()string{
	return fmt.Sprintf("%v: %v",pse.f,pse.e)
}

func (pse *PubSubError)Error()string{
	return fmt.Sprint(pse.e)
}

func (pse *PubSubError)Subscriber()interface{}{
	return pse.f
}

type wrap struct{
	f interface{}
}

func NewWrap(f interface{})*wrap{
	return &wrap{f:f}
}

type PubSub struct{
	c chan interface{}
	w []*wrap
	m sync.Mutex
	e chan error
}

func New() *PubSub{
	ps := new(PubSub)
	ps.c = make(chan interface{})
	ps.e = make(chan error)
	call := func(f interface{},rf reflect.Value,in []reflect.Value){
		defer func(){
			if err := recover();err != nil{
				ps.e <- &PubSubError{f,err}
			}
		}()
		rf.Call(in)
	}

	go func(){
		for v := range ps.c{
			rv := reflect.ValueOf(v)
			ps.m.Lock()
			for _,w := range ps.w{
				rf := reflect.ValueOf(w.f)
				if rv.Type() == reflect.ValueOf(w.f).Type().In(0){
					go call(w.f,rf,[]reflect.Value{rv})
				}
			}
			ps.m.Unlock()
		}
	}()
	return ps
}

func (ps *PubSub)Error()chan error{
	return ps.e
}

func (ps *PubSub)Sub(f interface{})error{
	check := f
	w,wrapped := f.(*wrap)
	if wrapped{
		check = w.f
	}

	rf := reflect.ValueOf(check)
	if rf.Kind() != reflect.Func{
		return errors.New("Not a function")
	}
	if rf.Type().NumIn() != 1{
		return errors.New("Number of arguments shoule be 1")
	}
	ps.m.Lock()
	defer ps.m.Unlock()

	if w,wrapped := f.(*wrap);wrapped{
		ps.w = append(ps.w,w)
	}else{
		ps.w = append(ps.w,&wrap{f:f})
	}

	return nil
}

func (ps *PubSub) Leave(f interface{}){
	var fp uintptr
	if f == nil{
		if pc,_,_,ok := runtime.Caller(1);ok{
			fp = runtime.FuncForPC(pc).Entry()
		}
	}else{
		fp = reflect.ValueOf(f).Pointer()
	}
	ps.m.Lock()
	defer ps.m.Unlock()

	result := make([]*wrap,0,len(ps.w))
	last := 0
	for i,v := range ps.w{
		vf := v.f
		if reflect.ValueOf(vf).Pointer() == fp{
			result = append(result,ps.w[last:i]...)
			last = i + 1
		}
	}
	ps.w = append(result,ps.w[last:]...)
}

func (ps *PubSub)Pub(v interface{}){
	ps.c <- v
}

func (ps *PubSub)Close(){
	close(ps.c)
	ps.w = nil
}