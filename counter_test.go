package counters

import (
	"fmt"
	"testing"
)

func TestWriteTo(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	box.GetCounter("test").Increment()
	fmt.Println(box.String())
}

func TestIncrement(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	cnt := box.GetCounter("test")
	cnt.Increment()
	cnt.IncrementBy(7)

	if v := box.GetCounter("test"); v.Value() != 8 {
		t.Errorf("got %d, expected 8", v.Value())
	}
}

func TestIncrementParallel(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	end := make(chan bool, 10)
	for x := 0; x < 10; x++ {
		go func() {
			for y := 0; y < 100; y++ {
				cnt := box.GetCounter("test")
				cnt.Increment()
				cnt.IncrementBy(3)
			}
			end <- true
		}()
	}
	for i := 0; i < 10; {
		if _, ok := <-end; ok {
			i++
		}
	}

	if v := box.GetCounter("test"); v.Value() != 4000 {
		t.Errorf("got %d, expected 4000", v.Value())
	}
}

func TestMax(t *testing.T) {
	box := NewCounterBox()
	r := box.GetMax("Olsztyn")
	r.Set(5)
	r.Set(10)
	r.Set(7)
	if v := box.GetMax("Olsztyn").Value(); v != 10 {
		t.Errorf("Max, want: 10, got %d", v)
	}
}

func TestPrefix(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	pref := box.WithPrefix("prefix:")
	cnt := pref.GetCounter("test")
	cnt.Increment()
	cnt.IncrementBy(7)

	if v := box.GetCounter("prefix:test"); v.Value() != 8 {
		t.Errorf("got %d, expected 8", v.Value())
	}
	if v := pref.GetCounter("test"); v.Value() != 8 {
		t.Errorf("got %d, expected 8", v.Value())
	}
}

func BenchmarkCounters(b *testing.B) {
	b.StopTimer()
	e := make(chan bool)
	c := NewCounterBox()
	f := func(b *testing.B, c *CounterBox, e chan bool) {
		for i := 0; i < b.N; i++ {
			c.GetCounter("abc123").IncrementBy(5)
			c.GetCounter("def456").IncrementBy(5)
			c.GetCounter("ghi789").IncrementBy(5)
			c.GetCounter("abc123").IncrementBy(5)
			c.GetCounter("def456").IncrementBy(5)
			c.GetCounter("ghi789").IncrementBy(5)
		}
		e <- true
	}
	b.StartTimer()
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)

	<-e
	<-e
	<-e
	<-e
	<-e
}

func BenchmarkCountersCached(b *testing.B) {
	b.StopTimer()
	e := make(chan bool)
	c := NewCounterBox()
	f := func(b *testing.B, c *CounterBox, e chan bool) {
		x := c.GetCounter("abc123")
		y := c.GetCounter("def456")
		z := c.GetCounter("ghi789")
		for i := 0; i < b.N; i++ {
			x.IncrementBy(5)
			y.IncrementBy(5)
			z.IncrementBy(5)
			x.IncrementBy(5)
			y.IncrementBy(5)
			z.IncrementBy(5)
		}
		e <- true
	}
	b.StartTimer()
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)
	go f(b, c, e)

	<-e
	<-e
	<-e
	<-e
	<-e
}
