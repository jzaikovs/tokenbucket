package tokenbucket

import "time"

// Bucket represents single token bucket
type Bucket struct {
	bucket
	capacity float32 // max tokens in bucket
	fillRate float32 // tokens per seconds
}

// New is constructor for creating single tokenbucket
func New(capacity, refillTime float32) *Bucket {
	return &Bucket{bucket: bucket{last: time.Now(), freeSpace: capacity}, capacity: capacity, fillRate: capacity / refillTime}
}

// Add adds token into bucket, returns true if token added or false if bucket is full
func (b *Bucket) Add(t time.Time) bool {
	return b.fill(b.capacity, b.fillRate, t)
}

type bucket struct {
	last      time.Time // last time bucket was used
	freeSpace float32   // free space in bucket
}

func (r *bucket) fill(capacity, fillrate float32, t time.Time) bool {
	r.freeSpace += float32(t.Sub(r.last).Seconds()) * fillrate
	r.last = t

	if r.freeSpace > capacity {
		r.freeSpace = capacity
	}

	if r.freeSpace < 1.0 {
		return false
	}

	r.freeSpace -= 1.0
	return true
}
