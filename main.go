package main

import (
	"runtime"
	"time"
	"reflect"
	"sync"
	"fmt"
)

var wait sync.WaitGroup

type MultiChannel struct{
	c chan interface{}
	f []interface{}
}

func New() *MultiChannel{
	mc := new(MultiChannel)
	mc.c = make(chan interface{})

	go func(){
		for v := range mc.c{
			rv := reflect.ValueOf(v)
			for _,f := range mc.f{
				rf := reflect.ValueOf(f)
				if rv.Type() == reflect.ValueOf(f).Type().In(0){
					rf.Call([]reflect.Value{rv})
				}
			}
		}
	}()

	return mc
}

func (mc *MultiChannel)Sub(f interface{})error{
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func{

	}

	if rf.Type().NumIn() != 1{

	}

	mc.f = append(mc.f,f)
	return nil
}

func (mc *MultiChannel)Pub(v interface{}){
	mc.c <-v
}

func (mc *MultiChannel)Leave(f interface{}){
	var fp uintptr

	if f == nil{
		if pc,_,_,ok := runtime.Caller(1);ok{
			fp = runtime.FuncForPC(pc).Entry()
		}
	}else{
		fp = reflect.ValueOf(f).Pointer()
	}

	result := make([]interface{},0,len(mc.f))
	last := 0

	for i,v := range mc.f{
		if reflect.ValueOf(v).Pointer() == fp{
			result = append(result,mc.f[last:i]...)
			last = i + 1
		}
	}
	mc.f = append(result,mc.f[last:]...)
}


func main1(){
	mc := New()
	mc.Sub(func(i int){
		fmt.Println("int subscriber: ",i)
	})
	mc.Sub(func(s string){
		fmt.Println("string subscriber: ",s)
	})

	mc.Pub(1)
	mc.Pub("hello")
	mc.Pub(2)

	time.Sleep(5 * time.Second)
}


func test(){
	test2()
}

func test2(){
	pc,file,line,ok := runtime.Caller(0)
	fmt.Println(pc)
	fmt.Println(file)
	fmt.Println(line)
	fmt.Println(ok)

	f := runtime.FuncForPC(pc)
	fmt.Println(f.Name())

	pc,file,line,ok = runtime.Caller(1)
	fmt.Println(pc)
	fmt.Println(file)
	fmt.Println(line)
	fmt.Println(ok)

	f = runtime.FuncForPC(pc)
	fmt.Println(f.Name())

	pc,file,line,ok = runtime.Caller(2)
	fmt.Println(pc)
	fmt.Println(file)
	fmt.Println(line)
	fmt.Println(ok)

	f = runtime.FuncForPC(pc)
	fmt.Println(f.Name())
}

func main2(){
	test()
}

func main(){
	fmt.Println("hello gdb")
}