package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cnt "github.com/orian/counters/global"
)

func main() {
	FLAG_signal := flag.Bool("signal", false, "handle SIGINT and SIGTERM")
	flag.Parse()
	if *FLAG_signal {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			lastInt := time.Now()
			for sig := range sigs {
				fmt.Printf("Got signal: %s(%d)\n", sig, sig)
				fmt.Printf("I am: %d\n", os.Getpid())
				fmt.Println(cnt.String())
				l := time.Now()
				if sig == syscall.SIGTERM || l.Sub(lastInt).Seconds() < 1. {
					os.Exit(0)
				}
				lastInt = l
			}
		}()
	}

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
