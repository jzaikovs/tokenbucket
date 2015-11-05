package tokenbucket

import (
	"sync"
	"time"
)

// Buckets is structure for handling spam-wall filtering
// It uses ratelimit tocket-bucket algorithm for each checked word
type Buckets struct {
	lock     sync.Mutex
	history  map[string]bucket
	capacity float32 // max tokens in bucket
	fillRate float32 // tokens per seconds
	gcSleep  time.Duration
}

// NewBuckets is constructor for buckets handler
func NewBuckets(capacity, refilTime float32, gcTime time.Duration) *Buckets {
	if capacity < 1.0 {
		panic("Capacity should be larger than 1.0")
	}

	buckets := &Buckets{history: make(map[string]bucket), capacity: capacity, fillRate: capacity / refilTime, gcSleep: gcTime}
	go buckets.gc()
	return buckets
}

// GC deletes old buckets from memory
func (buckets *Buckets) gc() {
	for {
		time.Sleep(buckets.gcSleep)
		buckets.lock.Lock()
		now := time.Now()
		for k, r := range buckets.history {
			// if bucket is in memory for too long then remove it
			// anyways on next use it will be empty so no reason to hold it in memory
			if float32(now.Sub(r.last).Seconds())*buckets.fillRate > buckets.capacity {
				delete(buckets.history, k)
			}
		}
		buckets.lock.Unlock()
	}
}

// Add adds token in bucket with specified key
func (buckets *Buckets) Add(key string) (ok bool) {
	now := time.Now()

	buckets.lock.Lock()
	defer buckets.lock.Unlock()

	r, ok := buckets.history[key]
	if !ok {
		// first occurrence
		buckets.history[key] = bucket{last: time.Now(), freeSpace: buckets.capacity - 1.0}
		return true
	}
	ok = (&r).fill(buckets.capacity, buckets.fillRate, now)

	buckets.history[key] = r

	return ok
}
