package counter

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

type MaxMinValue interface {
	Set(int)
	Name() string
	Value() int64
}

type Counter interface {
	Increment()
	IncrementBy(num int)
	Name() string
	Value() int64
}

type CounterBox struct {
	counters map[string]*counterImpl
	min      map[string]*minImpl
	max      map[string]*maxImpl
	m        *sync.RWMutex
}

func NewCounterBox() *CounterBox {
	return &CounterBox{
		counters: make(map[string]*counterImpl),
		min:      make(map[string]*minImpl),
		max:      make(map[string]*maxImpl),
		m:        &sync.RWMutex{},
	}
}

func (c *CounterBox) CreateHttpHandler() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.m.RLock()
		defer c.m.RUnlock()
		for k, v := range c.counters {
			fmt.Fprintf(w, "%s = %d\n", k, v.Value())
		}
	}
}

func (c *CounterBox) GetCounter(name string) Counter {
	c.m.RLock()
	if v, ok := c.counters[name]; ok {
		c.m.RUnlock()
		return v
	}
	c.m.RUnlock()
	c.m.Lock()
	defer c.m.Unlock()

	v := &counterImpl{name, 0}
	c.counters[name] = v
	return v
}

func (c *CounterBox) GetMin(name string) MaxMinValue {
	c.m.RLock()
	if v, ok := c.min[name]; ok {
		c.m.RUnlock()
		return v
	}
	c.m.RUnlock()
	c.m.Lock()
	defer c.m.Unlock()

	v := &minImpl{name, 0}
	c.min[name] = v
	return v
}

func (c *CounterBox) GetMax(name string) MaxMinValue {
	c.m.RLock()
	if v, ok := c.max[name]; ok {
		c.m.RUnlock()
		return v
	}
	c.m.RUnlock()
	c.m.Lock()
	defer c.m.Unlock()

	v := &maxImpl{name, 0}
	c.max[name] = v
	return v
}

type counterImpl struct {
	name  string
	value int64
}

func (c *counterImpl) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *counterImpl) IncrementBy(num int) {
	atomic.AddInt64(&c.value, int64(num))
}

func (c *counterImpl) Name() string {
	return c.Name()
}

func (c *counterImpl) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

type maxImpl counterImpl

func (m *maxImpl) Set(v int) {
	done := false
	v64 := int64(v)
	for !done {
		if o := atomic.LoadInt64(&m.value); v64 > o {
			done = atomic.CompareAndSwapInt64(&m.value, o, v64)
		} else {
			done = true
		}
	}
}

func (m *maxImpl) Name() string {
	return m.Name()
}

func (m *maxImpl) Value() int64 {
	return atomic.LoadInt64(&m.value)
}

type minImpl counterImpl

func (m *minImpl) Set(v int) {
	done := false
	v64 := int64(v)
	for !done {
		if o := atomic.LoadInt64(&m.value); v64 < o {
			done = atomic.CompareAndSwapInt64(&m.value, o, v64)
		} else {
			done = true
		}
	}
}

func (m *minImpl) Name() string {
	return m.Name()
}

func (m *minImpl) Value() int64 {
	return atomic.LoadInt64(&m.value)
}
