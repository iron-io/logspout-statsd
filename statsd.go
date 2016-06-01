package statsd

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"github.com/gliderlabs/logspout/router"
	"github.com/kr/logfmt"
	"github.com/quipo/statsd"
)

func init() {
	router.AdapterFactories.Register(NewStatsdAdapter, "statsd")
}

func NewStatsdAdapter(route *router.Route) (router.LogAdapter, error) {

	prefix := ""
	statsdclient := statsd.NewStatsdClient(route.Address, prefix)
	err := statsdclient.CreateSocket()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	// defer stats.Close()

	return &StatsdAdapter{
		route:        route,
		statsdClient: statsdclient,
		statsBuffer:  stats,
	}, nil
}

type StatsdAdapter struct {
	route        *router.Route
	counter      uint64
	statsdClient *statsd.StatsdClient
	statsBuffer  *statsd.StatsdBuffer
}

type Metric struct {
	Metric string // Name of the metric
	Value  int64
	Type   string
}

func (a *StatsdAdapter) Stream(logstream chan *router.Message) {
	for message := range logstream {
		if message.Container.Name == "/logspout" {
			continue
		}
		atomic.AddUint64(&a.counter, 1)
		// log.Println(atomic.LoadUint64(&a.counter), "source:", message.Source, "cname:", message.Container.Name, "mdata:", message.Data)

		m := &Metric{}
		if err := logfmt.Unmarshal([]byte(message.Data), m); err != nil {
			// log.Println("not in logfmt format, skipping")
			continue
		}
		// log.Println("metric:", *m)
		if m.Metric == "" {
			// log.Println("not a metric, skipping")
			continue
		}
		if m.Metric != "" {
			switch m.Type {
			case "count":
				a.statsBuffer.Incr(m.Metric, m.Value)
			case "gauge":
				a.statsBuffer.Gauge(m.Metric, m.Value)
			}
		}

	}
}
