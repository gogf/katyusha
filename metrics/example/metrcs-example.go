package main

import (
	"errors"
	"github.com/gogf/katyusha/metrics"
	"math/rand"
	"time"
)

const fanout = 10

func main() {
	c := metrics.Counter("foo")
	for i := 0; i < fanout; i++ {
		go func() {
			for {

				c.Dec(19)
				time.Sleep(time.Millisecond * 500)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(time.Millisecond * 400)
			}
		}()
	}

	g := metrics.Gauge("bar")
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(time.Millisecond * 300)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(time.Millisecond * 400)
			}
		}()
	}

	gf := metrics.GaugeFloat64("barfloat64")
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				gf.Update(19.0)
				time.Sleep(time.Millisecond * 300)
			}
		}()
		go func() {
			for {
				gf.Update(47.0)
				time.Sleep(time.Millisecond * 400)
			}
		}()
	}

	hc := metrics.NewHealthcheck(func(h metrics.IHealthcheck) {
		if 0 < rand.Intn(2) {
			h.Healthy()
		} else {
			h.Unhealthy(errors.New("baz"))
		}
	})
	_ = metrics.Register("baz", hc)

	s := metrics.NewExpDecaySample(1028, 0.015)
	h := metrics.Histogram("bang", s)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(time.Millisecond * 300)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(time.Millisecond * 400)
			}
		}()
	}

	m := metrics.Meter("guux")
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				m.Mark(19)

				time.Sleep(time.Millisecond * 300)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(time.Millisecond * 400)
			}
		}()
	}

	t := metrics.Timer("hooah")
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				t.Time(func() { time.Sleep(time.Millisecond * 300) })
			}
		}()
		go func() {
			for {
				t.Time(func() { time.Sleep(time.Millisecond * 400) })
			}
		}()
	}
	//
	//metrics.RegisterDebugGCStats(r)
	//go metrics.CaptureDebugGCStats(r, time.Second)
	//
	//metrics.RegisterRuntimeMemStats(r)
	//go metrics.CaptureRuntimeMemStats(r, time.Second)

	metrics.Log(time.Second)

	/*
		w, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrix")
		if nil != err { log.Fatalln(err) }
		metrix.Syslog(r, 60e9, w)
	*/

	/*
		addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:2003")
		metrix.Graphite(r, 10e9, "metrix", addr)
	*/

	/*
		stathat.Stathat(r, 10e9, "example@example.com")
	*/

}
