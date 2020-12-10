package main

import (
	"errors"
	"github.com/gogf/katyusha/metrix"
	"math/rand"
	"time"
)

const fanout = 10

func main() {
	r := metrix.NewRegistry()

	c := metrix.NewCounter()
	r.Register("foo", c)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				c.Dec(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(400e6)
			}
		}()
	}

	g := metrix.NewGauge()
	r.Register("bar", g)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	gf := metrix.NewGaugeFloat64()
	r.Register("barfloat64", gf)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19.0)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47.0)
				time.Sleep(400e6)
			}
		}()
	}

	hc := metrix.NewHealthcheck(func(h metrix.IHealthcheck) {
		if 0 < rand.Intn(2) {
			h.Healthy()
		} else {
			h.Unhealthy(errors.New("baz"))
		}
	})
	r.Register("baz", hc)

	s := metrix.NewExpDecaySample(1028, 0.015)
	//s := metrix.NewUniformSample(1028)
	h := metrix.NewHistogram(s)
	r.Register("bang", h)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	m := metrix.NewMeter()
	r.Register("quux", m)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				m.Mark(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(400e6)
			}
		}()
	}

	t := metrix.NewTimer()
	r.Register("hooah", t)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				t.Time(func() { time.Sleep(300e6) })
			}
		}()
		go func() {
			for {
				t.Time(func() { time.Sleep(400e6) })
			}
		}()
	}

	metrix.RegisterDebugGCStats(r)
	go metrix.CaptureDebugGCStats(r, 5e9)

	metrix.RegisterRuntimeMemStats(r)
	go metrix.CaptureRuntimeMemStats(r, 5e9)

	metrix.Log(r, 60e9)

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
