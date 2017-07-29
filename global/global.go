// Package global provides a package level access to all
// counters.CounterBox methods. Effectively it's a global variable
// for counters.

package global

import (
	"io"
	"net/http"
	"time"

	"github.com/orian/counters"
	"github.com/sirupsen/logrus"
)

var globalBox *counters.CounterBox

func init() {
	globalBox = counters.NewCounterBox()
}

func GetCounter(name string) counters.Counter {
	return globalBox.GetCounter(name)
}

func Get(name string) counters.Counter {
	return globalBox.GetCounter(name)
}

func Min(name string) counters.MaxMinValue {
	return globalBox.GetMin(name)
}

func Max(name string) counters.MaxMinValue {
	return globalBox.GetMax(name)
}

func WithPrefix(prefix string) counters.Counters {
	return globalBox.WithPrefix(prefix)
}

func GetMin(name string) counters.MaxMinValue {
	return globalBox.GetMin(name)
}

func GetMax(name string) counters.MaxMinValue {
	return globalBox.GetMax(name)
}

func CreateHttpHandler() http.HandlerFunc {
	return globalBox.CreateHttpHandler()
}

func Default() counters.Counters {
	return globalBox
}

func WriteTo(w io.Writer) {
	globalBox.WriteTo(w)
}

func String() string {
	return globalBox.String()
}

func LogrusOnSignal() {
	counters.InitCountersOnSignal(logrus.StandardLogger(), globalBox)
}

func LogrusCountersEvery(d time.Duration) {
	counters.LogCountersEvery(logrus.StandardLogger(), globalBox, d)
}
