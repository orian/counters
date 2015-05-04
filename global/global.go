package global

import (
	"github.com/orian/counters"
)

var counterBox *counters.CounterBox

func init() {
	counterBox = counters.NewCounterBox()
}

func GetCounter(name string) counters.Counter {
	return counterBox.GetCounter(name)
}

func GetMin(name string) MaxMinValue {
	return counterBox.GetMin(name)
}

func GetMax(name string) MaxMinValue {
	return counterBox.GetMax(name)
}

func CreateHttpHandler() HandlerFunc {
	return counterBox.CreateHttpHandler()
}
