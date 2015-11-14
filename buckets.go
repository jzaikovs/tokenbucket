package tokenbucket

import (
	"sync"
	"time"
)

// Buckets is structure for handling multiple token buckets with same rate and initial capacity
type Buckets struct {
	lock     sync.Mutex
	buckets  map[string]bucket
	capacity float32 // max tokens in bucket
	rate     float32 // tokens per seconds
}

// NewBuckets is constructor for buckets handler, with initial bucket capacity and rate (tokens/second) at which buckets gain free space
func NewBuckets(capacity int, rate float32) *Buckets {
	if capacity < 1 {
		panic("Capacity should be larger than 1")
	}

	buckets := &Buckets{buckets: make(map[string]bucket), capacity: float32(capacity), rate: rate}

	// start bucket cleanup with interval in witch buckets can fully refill if not used
	go buckets.freeBuckets(time.Duration(float32(capacity)/rate)*time.Second + time.Second)

	return buckets
}

// freeBuckets deletes old buckets from memory and sleeps for some duration before cleans again
func (buckets *Buckets) freeBuckets(sleep time.Duration) {
	for {
		time.Sleep(sleep)
		buckets.lock.Lock()
		now := time.Now()
		for k, r := range buckets.buckets {
			// if bucket is in memory for too long then remove it
			// anyways on next use it will be empty so no reason to hold it in memory
			if float32(now.Sub(r.last).Seconds())*buckets.rate > buckets.capacity {
				delete(buckets.buckets, k)
			}
		}
		buckets.lock.Unlock()
	}
}

// Add adds token in bucket with specified name and returns free space left in bucket and ok if token was added
func (buckets *Buckets) Add(name string, t time.Time) (space int, ok bool) {
	buckets.lock.Lock()
	defer buckets.lock.Unlock()

	r, ok := buckets.buckets[name]
	if !ok {
		// first occurrence
		b := bucket{last: t, space: buckets.capacity - 1.0}
		buckets.buckets[name] = b
		return int(b.space), true
	}

	space, ok = (&r).fill(buckets.capacity, buckets.rate, t)

	buckets.buckets[name] = r

	return space, ok
}

// Capacity returns maximum bucket capacity
func (buckets *Buckets) Capacity() int {
	return int(buckets.capacity)
}

// Check adds token in specidfied bucket and returns true if token added
func (buckets *Buckets) Check(name string, t time.Time) bool {
	_, ok := buckets.Add(name, t)
	return ok
}
