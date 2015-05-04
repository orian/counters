package main

import (
	cnt "github.com/orian/counters/global"
	"net/http"
	"time"
)

func main() {
	cnt.GetCounter("start").Increment()
	http.Handle("/status", cnt.CreateHttpHandler())
	go func() {
		c := time.Tick(1 * time.Second)
		for range c {
			cnt.GetCounter("ticker").Increment()
		}
	}()
	cnt.GetMax("monist").Set(128)
	http.ListenAndServe(":8080", nil)
}
