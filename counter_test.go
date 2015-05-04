package counter

import (
	"testing"
)

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
