package main

import (
	"github.com/gogf/katyusha/metrics"
	"time"
)

func main() {
	t := metrics.Timer("hooah")
	t.Time(func() {
		time.Sleep(time.Second)
	})

	go metrics.Log(time.Second)

	var j int64
	for j < 30 {
		time.Sleep(time.Second)
		j++
		//h.Update(j)
	}
}
