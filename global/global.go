package global

import (
	"github.com/orian/counter"
)

var counterBox *counter.CounterBox

func init() {
	counterBox = counter.NewCounterBox()
}

func GetCounter(name string) counter.Counter {
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
