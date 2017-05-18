package mnemosyned

import (
	"sync"
	"time"

	"github.com/piotrkowalczuk/mnemosyne/mnemosynerpc"
)

type cacheEntry struct {
	ses     mnemosynerpc.Session
	exp     time.Time
	refresh bool
}

type cache struct {
	monitor  monitoringCache
	data     map[uint64]*cacheEntry
	dataLock sync.RWMutex
	ttl      time.Duration
}

func (c *cache) refresh(k uint64) {
	if c.monitor.enabled {
		c.monitor.refresh.Add(1)
	}
	c.dataLock.Lock()
	c.data[k].refresh = true
	c.dataLock.Unlock()
}

func (c *cache) put(k uint64, ses mnemosynerpc.Session) {
	c.dataLock.Lock()
	c.data[k] = &cacheEntry{ses: ses, exp: time.Now().Add(c.ttl), refresh: false}
	c.dataLock.Unlock()
}

func (c *cache) del(k uint64) {
	c.dataLock.Lock()
	delete(c.data, k)
	c.dataLock.Unlock()
}

func (c *cache) read(k uint64) (*cacheEntry, bool) {
	c.dataLock.RLock()
	entry, ok := c.data[k]
	c.dataLock.RUnlock()
	if c.monitor.enabled {
		if ok {
			c.monitor.hits.Add(1)
		} else {
			c.monitor.misses.Add(1)
		}
	}
	return entry, ok
}
