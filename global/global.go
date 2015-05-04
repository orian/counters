package global

import (
	"github.com/orian/counter"
)

var counterBox = counter.NewCounterBox()

func init() {
	counterBox.Start()
}

func GetCounter(name string) counter.Counter {
	return counterBox.GetCounter(name)
}
