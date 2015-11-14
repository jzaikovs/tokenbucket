package tokenbucket

import "time"

// Bucket represents single token bucket
type Bucket struct {
	bucket
	capacity float32 // max tokens in bucket
	rate     float32 // tokens per seconds
}

// New is constructor for creating single token bucket
// with specific capacity (token count) and rate (tokens/second) at which buckets gain free space
func New(capacity int, rate float32) *Bucket {
	if capacity < 1 {
		panic("Capacity should be larger than 1")
	}

	return &Bucket{bucket: bucket{last: time.Now(), space: float32(capacity)}, capacity: float32(capacity), rate: rate}
}

// Add adds token into bucket, returns free space left in bucket and ok if token was added
func (b *Bucket) Add(t time.Time) (int, bool) {
	return b.fill(b.capacity, b.rate, t)
}

// Check adds token in bucket and returns true if token added
func (b *Bucket) Check(t time.Time) bool {
	_, ok := b.fill(b.capacity, b.rate, t)
	return ok
}

type bucket struct {
	last  time.Time // last time bucket was used
	space float32   // free space (in tokens) in bucket
}

func (r *bucket) fill(capacity, rate float32, t time.Time) (space int, ok bool) {
	r.space += float32(t.Sub(r.last).Seconds()) * rate
	r.last = t

	if r.space > capacity {
		r.space = capacity
	}

	if r.space < 1.0 {
		// bucket is full no more free space, so next token can be added only when there is space for token
		return int(r.space), false
	}

	r.space -= 1.0
	return int(r.space), true
}
