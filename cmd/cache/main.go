package main

import (
	"fmt"
	"github.com/b2wdigital/restQL-golang/internal/platform/cache"
	"sync/atomic"
	"time"
)

var counter uint64 = 0

func calcValue(key interface{}, value interface{}) (interface{}, error) {
	fmt.Printf("load executed : key=%s\n", key)
	v := fmt.Sprintf("ok-%d", counter)
	atomic.AddUint64(&counter, 1)

	return v, nil
}

func main() {
	c := cache.New(nil, calcValue, 20, cache.WithRefreshInterval(50*time.Millisecond), cache.WithRefreshQueueLength(10))

	for j := 0; j < 50; j++ {
		for i := 0; i < 4; i++ {
			key := fmt.Sprintf("k%d", i)
			item, _ := c.Get(nil, key)
			fmt.Printf("cache value : %s=%v\n", key, item)
			time.Sleep(10 * time.Millisecond)
		}
	}

	item, _ := c.Get(nil, "k0")
	fmt.Printf("value of first key : %v", item)
}
