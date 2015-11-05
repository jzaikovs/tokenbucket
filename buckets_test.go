package tokenbucket

import (
	"testing"
	"time"
)

func TestFill(t *testing.T) {
	rate := 3
	buckets := NewBuckets(float32(rate), 1, time.Second/10)
	// using frequent GC collection, GC should not remove our bucket
	// this means that in bucket each second we can add 2 tokens
	for i := 0; i < 4; i++ {
		for j := 0; j < rate; j++ {
			if !buckets.Add("hello") {
				t.Error("wall blocked too fast", i, j)
				return
			}
		}
		if buckets.Add("hello") { // <- this should fail, bucket is full
			t.Error("wall didn't blocked at wall boundry")
			return
		}
		time.Sleep(time.Second) // wait for bucket to empty
	}
}

func TestGC(t *testing.T) {
	buckets := NewBuckets(10, 0.5, time.Second)
	buckets.Add("hello")
	time.Sleep(time.Second + time.Millisecond*100)
	if len(buckets.history) != 0 {
		t.Error("GC is not working")
	}
}

func TestBucket(t *testing.T) {
	rate := 3
	bucket := New(float32(rate), 1)
	// using frequent GC collection, GC should not remove our bucket
	// this means that in bucket each second we can add 2 tokens
	for i := 0; i < 4; i++ {
		for j := 0; j < rate; j++ {
			if !bucket.Add(time.Now()) {
				t.Error("wall blocked too fast", i, j)
				return
			}
		}
		if bucket.Add(time.Now()) { // <- this should fail, bucket is full
			t.Error("wall didn't blocked at wall boundry")
			return
		}
		time.Sleep(time.Second) // wait for bucket to empty
	}
}
