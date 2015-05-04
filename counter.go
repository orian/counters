package counter

import (
	"fmt"
)

type Counter interface {
	Increment()
	IncrementBy(num int)
}

type CounterBox struct {
	counters  map[string]int
	changes   chan update
	pleaseEnd chan bool
	done      chan bool
}

func NewCounterBox() *CounterBox {
	return &CounterBox{make(map[string]int), make(chan update, 10), make(chan bool, 1), make(chan bool, 1)}
}

func (box *CounterBox) GetCounter(name string) Counter {
	return &counterImpl{name, box.changes}
}

func (box *CounterBox) Start() {
	go func() {
		for {
			select {
			case change := <-box.changes:
				box.counters[change.name] += change.num
			case _, ok := <-box.pleaseEnd:
				if ok {
					box.done <- true
				} else {
					fmt.Printf("asking for end: !ok \n")
					return
				}
			}
		}
	}()
}

func (box *CounterBox) EndAsync() {
	close(box.changes)
}

func (box *CounterBox) IsDone() {
	<-box.done
}

func (box *CounterBox) End() {
	box.pleaseEnd <- true
	box.IsDone()
}

type update struct {
	name string
	num  int
}

type counterImpl struct {
	name string
	s    chan<- update
}

func (cnt *counterImpl) Increment() {
	cnt.s <- update{cnt.name, 1}
}

func (cnt *counterImpl) IncrementBy(num int) {
	cnt.s <- update{cnt.name, num}
}
