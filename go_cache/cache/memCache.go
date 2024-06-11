package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type memCache struct {
	maxMemorySize             int64
	maxMemorySizeStr          string
	currentMemorySize         int64
	values                    map[string]*memCacheValue
	lock                      sync.RWMutex
	cleanExpiredItemsInterval time.Duration
}
type memCacheValue struct {
	val        interface{}
	expireTime time.Time
	expire     time.Duration
	size       int64
}

func NewMemCache() Cache {
	mc := &memCache{
		values:                    make(map[string]*memCacheValue),
		cleanExpiredItemsInterval: time.Second,
	}
	go mc.cleanExpiredItem()
	return mc
}
func (mc *memCache) SetMaxMemory(size string) bool {
	mc.maxMemorySize, mc.maxMemorySizeStr = ParseSize(size)
	return true
}
func (mc *memCache) Set(key string, val interface{}, expire time.Duration) bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	v := &memCacheValue{
		val:        val,
		expireTime: time.Now().Add(expire),
		expire:     expire,
		size:       GetValSize(val),
	}
	if _, ok := mc.values[key]; ok {
		mc.del(key)
	}
	mc.add(key, v)

	if mc.currentMemorySize > mc.maxMemorySize {
		mc.del(key)
		log.Println(fmt.Sprintf("max memory size %s", mc.maxMemorySizeStr))
	}
	return false
}
func (mc *memCache) get(key string) (*memCacheValue, bool) {
	val, ok := mc.values[key]
	return val, ok
}
func (mc *memCache) del(key string) {
	tmp, ok := mc.values[key]
	if ok && tmp != nil {
		mc.currentMemorySize -= tmp.size
		delete(mc.values, key)
	}
}
func (mc *memCache) add(key string, val *memCacheValue) {
	mc.values[key] = val
	mc.currentMemorySize += val.size
}
func (mc *memCache) Get(key string) (interface{}, bool) {
	mc.lock.RLock()
	defer mc.lock.RUnlock()

	mcv, ok := mc.get(key)
	if ok {
		if mcv.expire != 0 && mcv.expireTime.Before(time.Now()) {
			mc.del(key)
			return nil, false
		}
		return mcv.val, ok
	}
	return nil, false
}
func (mc *memCache) Del(key string) bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.del(key)
	return true
}
func (mc *memCache) Exists(key string) bool {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	_, ok := mc.get(key)
	return ok
}
func (mc *memCache) Keys() int64 {
	mc.lock.RLock()
	defer mc.lock.RUnlock()
	return int64(len(mc.values))
}
func (mc *memCache) Flush() bool {
	mc.lock.Lock()
	defer mc.lock.Unlock()
	mc.values = make(map[string]*memCacheValue, 0)
	mc.currentMemorySize = 0
	return true
}
func (mc *memCache) cleanExpiredItem() {
	timeTicker := time.NewTicker(mc.cleanExpiredItemsInterval)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			for key, v := range mc.values {
				if v.expire != 0 && time.Now().After(v.expireTime) {
					mc.lock.Lock()
					mc.del(key)
					mc.lock.Unlock()
				}
			}
		}
	}
}
