package cachering

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	ttlRing := New(func(key string) interface{} {
		time.Sleep(100 * time.Millisecond)
		return fmt.Sprintf("%s_%d", key, time.Now().Nanosecond())
	}, 4)

	assert.NotNil(t, ttlRing, "can not initiate ring")

	for i := 0; i < 8; i++ {
		time.Sleep(25 * time.Millisecond)
		_ = ttlRing.Get("testkey_1", 75 * time.Millisecond)

	}

	stats := ttlRing.statsAgent.GetStats()

	fmt.Printf("Expired: %d\n", stats[EVENT_EXPIRED])
	fmt.Printf("Hit: %d\n", stats[EVENT_HIT])
	fmt.Printf("Miss: %d\n", stats[EVENT_MISS])
	fmt.Printf("Not found: %d\n", stats[EVENT_NOT_FOUND])

	assert.Less(t, stats[EVENT_MISS], stats[EVENT_HIT], "Miss count must be smaller than hit count")
	assert.Less(t, stats[EVENT_EXPIRED], stats[EVENT_HIT], "Expired count must be smaller than hit count")
}

func TestCacheRing(t *testing.T) {
	ttlRing := New(func(key string) interface{} {
		return fmt.Sprintf("%s_%d", key, time.Now().Nanosecond())
	}, 3)

	ttlRing.Get("testkey_1", 1 * time.Hour)
	ttlRing.Get("testkey_2", 1 * time.Hour)
	ttlRing.Get("testkey_3", 1 * time.Hour)

	ttlRing.Get("testkey_4", 1 * time.Hour)

	assert.Contains(t, ttlRing.store, "testkey_4")
	assert.NotContains(t, ttlRing.store, "testkey_1")

	ttlRing.Get("testkey_5", 1 * time.Hour)

	assert.Contains(t, ttlRing.store, "testkey_5")
	assert.NotContains(t, ttlRing.store, "testkey_2")


}