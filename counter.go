// Package counters provides a simple counter, max and min functionalities.
// All counters are kept in CounterBox.
// Library is thread safe.
package counters

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"text/template"
	"time"
)

// MaxMinValue is an interface for minima and maxima counters.
type MaxMinValue interface {
	// Set allows to update value if necessary.
	Set(int)
	// Name returns a name of counter.
	Name() string
	// Value returns a current value.
	Value() int64
}

// Counter is an interface for integer increase only counter.
type Counter interface {
	// Increment increases counter by one.
	Increment() int64
	// IncrementBy increases counter by a given number.
	IncrementBy(num int) int64
	// Decrement decreases counter by one.
	Decrement() int64
	// DecrementBy decreases counter by a given number.
	DecrementBy(num int) int64
	// Set sets a specific value.
	Set(num int)
	// Name returns a name of counter.
	Name() string
	// Value returns a current value of counter.
	Value() int64
}

type Counters interface {
	Get(string) Counter
	Min(string) MaxMinValue
	Max(string) MaxMinValue
	WithPrefix(string) Counters
	GetCounter(string) Counter
	GetMin(string) MaxMinValue
	GetMax(string) MaxMinValue
	WriteTo(w io.Writer)
	Prefix() string
	String() string
}

// CounterBox is a main type, it keeps references to all counters
// requested from it.
type CounterBox struct {
	counters *sync.Map
	min      *sync.Map
	max      *sync.Map
}

// NewCounterBox creates a new object to keep all counters.
func NewCounterBox() *CounterBox {
	return &CounterBox{
		counters: &sync.Map{},
		min:      &sync.Map{},
		max:      &sync.Map{},
	}
}

func New() Counters {
	return NewCounterBox()
}

// CreateHttpHandler creates a simple handler printing values of all counters.
func (c *CounterBox) CreateHttpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { c.WriteTo(w) }
}

func (c *CounterBox) Get(name string) Counter {
	return c.GetCounter(name)
}

func (c *CounterBox) Min(name string) MaxMinValue {
	return c.GetMin(name)
}

func (c *CounterBox) Max(name string) MaxMinValue {
	return c.GetMax(name)
}

type prefixed struct {
	CounterBox
	base   *CounterBox
	prefix string
}

func (c *CounterBox) WithPrefix(name string) Counters {
	return &prefixed{
		CounterBox{
			counters: &sync.Map{},
			min:      &sync.Map{},
			max:      &sync.Map{},
		},
		c,
		name}
}

func (c *CounterBox) Prefix() string {
	return ""
}

func (c *prefixed) GetCounter(name string) Counter {
	value, _ := c.counters.LoadOrStore(name, c.base.GetCounter(c.prefix+name))
	v, _ := value.(Counter)
	return v
}

// GetMin returns a minima counter of given name, if doesn't exist than create.
func (c *prefixed) GetMin(name string) MaxMinValue {
	value, _ := c.min.LoadOrStore(name, c.base.GetMin(c.prefix+name))
	v, _ := value.(MaxMinValue)
	return v
}

// GetMax returns a maxima counter of given name, if doesn't exist than create.
func (c *prefixed) GetMax(name string) MaxMinValue {
	value, _ := c.max.LoadOrStore(name, c.base.GetMax(c.prefix+name))
	v, _ := value.(MaxMinValue)
	return v
}

func (c *prefixed) Get(name string) Counter {
	return c.GetCounter(name)
}

// GetMin returns a minima counter of given name, if doesn't exist than create.
func (c *prefixed) Min(name string) MaxMinValue {
	return c.GetMin(name)
}

// GetMax returns a maxima counter of given name, if doesn't exist than create.
func (c *prefixed) Max(name string) MaxMinValue {
	return c.GetMax(name)
}

func (c *prefixed) Prefix() string {
	return c.prefix
}

// GetCounter returns a counter of given name, if doesn't exist than create.
func (c *CounterBox) GetCounter(name string) Counter {
	value, _ := c.counters.LoadOrStore(name, &counterImpl{name, 0})
	v, _ := value.(Counter)
	return v
}

// GetMin returns a minima counter of given name, if doesn't exist than create.
func (c *CounterBox) GetMin(name string) MaxMinValue {
	value, _ := c.min.LoadOrStore(name, &minImpl{name, math.MaxInt64})
	v, _ := value.(MaxMinValue)
	return v
}

// GetMax returns a maxima counter of given name, if doesn't exist than create.
func (c *CounterBox) GetMax(name string) MaxMinValue {
	value, _ := c.max.LoadOrStore(name, &maxImpl{name, 0})
	v, _ := value.(MaxMinValue)
	return v
}

var tmpl = template.Must(template.New("main").Parse(`== Counters ==
{{- range .Counters}}
  {{.Name}}: {{.Value}}
{{- end}}
== Min values ==
{{- range .Min}}
  {{.Name}}: {{.Value}}
{{- end}}
== Max values ==
{{- range .Max}}
  {{.Name}}: {{.Value}}
{{- end -}}
`))

func (c *CounterBox) WriteTo(w io.Writer) {
	data := &struct {
		Counters []Counter
		Min      []MaxMinValue
		Max      []MaxMinValue
	}{}
	c.counters.Range(func(key interface{}, value interface{}) bool {
		if value, ok := value.(Counter); ok {
			data.Counters = append(data.Counters, value)
		}
		return true
	})
	c.min.Range(func(key interface{}, value interface{}) bool {
		if value, ok := value.(MaxMinValue); ok {
			data.Min = append(data.Min, value)
		}
		return true
	})
	c.max.Range(func(key interface{}, value interface{}) bool {
		if value, ok := value.(MaxMinValue); ok {
			data.Max = append(data.Max, value)
		}
		return true
	})
	sort.Slice(data.Counters, func(i, j int) bool { return strings.Compare(data.Counters[i].Name(), data.Counters[j].Name()) < 0 })
	sort.Slice(data.Min, func(i, j int) bool { return strings.Compare(data.Min[i].Name(), data.Min[j].Name()) < 0 })
	sort.Slice(data.Max, func(i, j int) bool { return strings.Compare(data.Max[i].Name(), data.Max[j].Name()) < 0 })
	tmpl.Execute(w, data)
}

func (c *CounterBox) String() string {
	buf := &bytes.Buffer{}
	c.WriteTo(buf)
	return buf.String()
}

type counterImpl struct {
	name  string
	value int64
}

func (c *counterImpl) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *counterImpl) IncrementBy(num int) int64 {
	return atomic.AddInt64(&c.value, int64(num))
}

func (c *counterImpl) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

func (c *counterImpl) DecrementBy(num int) int64 {
	return atomic.AddInt64(&c.value, -int64(num))
}

func (c *counterImpl) Set(num int) {
	atomic.StoreInt64(&c.value, int64(num))
}

func (c *counterImpl) Name() string {
	return c.name
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
	return m.name
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
	return m.name
}

func (m *minImpl) Value() int64 {
	return atomic.LoadInt64(&m.value)
}

type TrivialLogger interface {
	Print(...interface{})
}

func InitCountersOnSignal(logger TrivialLogger, box Counters) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		lastInt := time.Now()
		for sig := range sigs {
			logger.Print(box.String())
			l := time.Now()
			if sig == syscall.SIGTERM || l.Sub(lastInt).Seconds() < 1. {
				os.Exit(0)
			}
			lastInt = l
		}
	}()
}

func LogCountersEvery(logger TrivialLogger, box Counters, d time.Duration) {
	go func() {
		t := time.NewTicker(d)
		for range t.C {
			logger.Print(box.String())
		}
	}()
}
