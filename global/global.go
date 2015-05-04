// Package global provides a package level access to all
// counters.CounterBox methods. Effectively it's a global variable
// for counters.

package global

import (
	"github.com/orian/counters"
	"net/http"
)

var counterBox *counters.CounterBox

func init() {
	counterBox = counters.NewCounterBox()
}

func GetCounter(name string) counters.Counter {
	return counterBox.GetCounter(name)
}

func GetMin(name string) counters.MaxMinValue {
	return counterBox.GetMin(name)
}

func GetMax(name string) counters.MaxMinValue {
	return counterBox.GetMax(name)
}

func CreateHttpHandler() http.HandlerFunc {
	return counterBox.CreateHttpHandler()
}
