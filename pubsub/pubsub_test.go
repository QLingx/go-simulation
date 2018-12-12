package pubsub

import (
	"testing"
)

func TestInt(t *testing.T){
	done := make(chan int)
	ps := New()
	ps.Sub(func(i int){
		done <- i
	})
	ps.Pub(1)
	i := <- done
	if i != 1{
		t.Fatalf("Excepted %v, but %d:",1,i)
	}
}

type F struct{
	m string
}

func TestStruct(t *testing.T){
	done := make(chan *F)
	ps := New()

	ps.Sub(func(f *F){
		done <- f
	})

	ps.Pub(&F{
		"hello world",
	})

	f := <-done

	if f.m != "hello world"{
		t.Fatalf("Excepted %v,but %s:","hello world",f.m)
	}
}

func TestOnly(t *testing.T){
	doneInt := make(chan int)
	doneF := make(chan *F)
	ps := New()
	ps.Sub(func(i int){
		doneInt <- i
	})

	ps.Sub(func(f *F){
		doneF <- f
	})

	ps.Pub(&F{
		"hello world",
	})

	ps.Pub(2)

	i := <-doneInt
	f := <-doneF

	if f.m != "hello world"{
		t.Fatalf("Excepted %v,but %s:","hello world",f.m)
	}

	if i != 2{
		t.Fatalf("Excepted %v,but %s:",2,f.m)
	}
}

func TestLeave(t *testing.T){
	done := make(chan int)
	ps := New()

	f := func(i int){
		done <- i
	}
	ps.Sub(f)
	ps.Sub(f)
	ps.Pub(1)
	i1 := <-done
	i2 := <-done

	if i1 != 1 || i2 != 1{
		t.Fatal("excepted multiple subscribers")
	}

	ps.Leave(f)
	ps.Pub(2)
	select{
	case <-done:
		t.Fatal("wtf")
	default:
	}

}