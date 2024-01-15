package cache

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

// DefaultSize determines how big cache should be at the beginning.
const DefaultSize = 100000

type Entry struct {
	Ses     mnemosynerpc.Session
	Exp     time.Time
	Refresh bool
}

type Cache struct {
	data     map[uint64]*Entry
	dataLock sync.RWMutex
	TTL      time.Duration
	// monitoring
	hitsTotal    prometheus.Counter
	missesTotal  prometheus.Counter
	refreshTotal prometheus.Counter
}

func New(ttl time.Duration, namespace string) *Cache {
	return &Cache{
		TTL:  ttl,
		data: make(map[uint64]*Entry, DefaultSize),
		hitsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "cache",
			Name:      "hits_total",
			Help:      "Total number of cache hits.",
		}),
		missesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "cache",
			Name:      "misses_total",
			Help:      "Total number of cache misses.",
		}),
		refreshTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "cache",
			Name:      "refresh_total",
			Help:      "Total number of times cache Refresh.",
		}),
	}
}

func (c *Cache) Refresh(k uint64) {
	c.refreshTotal.Add(1)

	c.dataLock.Lock()
	c.data[k].Refresh = true
	c.dataLock.Unlock()
}

func (c *Cache) Put(k uint64, ses mnemosynerpc.Session) {
	c.dataLock.Lock()
	c.data[k] = &Entry{Ses: ses, Exp: time.Now().Add(c.TTL), Refresh: false}
	c.dataLock.Unlock()
}

func (c *Cache) Del(k uint64) {
	c.dataLock.Lock()
	delete(c.data, k)
	c.dataLock.Unlock()
}

func (c *Cache) Read(k uint64) (*Entry, bool) {
	c.dataLock.RLock()
	entry, ok := c.data[k]
	c.dataLock.RUnlock()
	if ok {
		c.hitsTotal.Add(1)
	} else {
		c.missesTotal.Add(1)
	}
	return entry, ok
}

// Collect implements prometheus Collector interface.
func (c *Cache) Collect(in chan<- prometheus.Metric) {
	c.hitsTotal.Collect(in)
	c.refreshTotal.Collect(in)
	c.missesTotal.Collect(in)
}

// Describe implements prometheus Collector interface.
func (c *Cache) Describe(in chan<- *prometheus.Desc) {
	c.hitsTotal.Describe(in)
	c.refreshTotal.Describe(in)
	c.missesTotal.Describe(in)
}
