package tokenbucket

import (
	"testing"
	"time"
)

func TestFill(t *testing.T) {
	// TODO: test multiple buckets with different name
	rate := 3
	buckets := NewBuckets(rate, float32(rate))
	// with this capacity and rate, bucket will be emptied after each second
	// that means we can add capacity amount of tokens in second
	// and adding one more will return false
	// it can be repeated each second, because after second bucket will be empty
	for i := 0; i < 4; i++ {
		for j := 0; j < rate; j++ {
			if !buckets.Add("hello", time.Now()) {
				t.Error("wall blocked too fast", i, j)
				return
			}
		}
		if buckets.Add("hello", time.Now()) { // <- this should fail, bucket is full
			t.Error("wall didn't blocked at wall boundry")
			return
		}
		time.Sleep(time.Second) // wait for bucket to empty
	}
}

func TestGC(t *testing.T) {
	buckets := NewBuckets(2, 4)
	buckets.Add("hello", time.Now())
	time.Sleep(time.Second + time.Millisecond*600)
	if len(buckets.buckets) != 0 {
		t.Error("GC is not working")
	}
}

func TestBucket(t *testing.T) {
	capacity := 3
	bucket := New(capacity, float32(capacity))
	// with this capacity and rate, bucket will be emptied after each second
	// that means we can add capacity amount of tokens in second
	// and adding one more will return false
	// it can be repeated each second, because after second bucket will be empty
	for i := 0; i < 4; i++ {
		for j := 0; j < capacity; j++ {
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
