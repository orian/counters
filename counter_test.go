package counter

import (
	"testing"
)

func TestIncrement(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	box.Start()
	cnt := box.GetCounter("test")
	cnt.Increment()
	cnt.IncrementBy(7)
	box.End()

	if v := box.counters["test"]; v != 8 {
		t.Errorf("got %d, expected 8", v)
	}
}

func TestIncrementParallel(t *testing.T) {
	t.Parallel()
	box := NewCounterBox()
	box.Start()
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
	box.End()

	if v := box.counters["test"]; v != 4000 {
		t.Errorf("got %d, expected 4000", v)
	}
}
