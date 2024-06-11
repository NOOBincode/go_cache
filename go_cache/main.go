package main

import (
	"go_cache/cache_server"
)

func main() {
	cache := cache_server.NewMemCache()
	cache.SetMaxMemory("100MB")
	cache.Set("int", 1)
	cache.Set("bool", false)
	cache.Set("data", map[string]interface{}{"a": 1})
	cache.Get("int")
	cache.Del("int")
	cache.Flush()
	cache.Keys()
}
