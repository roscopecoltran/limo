package service

import (
	"github.com/roscopecoltran/feedify/memcache"
)

func NewCache() *memcache.MemcacheClient {
	return memcache.NewMemcacheClient()
}
