package lru

import (
	"container/list"
	"go/types"
)

type Cache struct {
	maxBytes int64
	nBytes int64
	list	*list.List
	cache map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	Value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		list: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToBack(ele)
		kv := ele.Value.(*entry)
		return kv.Value, true
	}

	return
}

func (c *Cache) Remove() {
	ele := c.list.Front()

	if ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.Value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.Value)
		}
	}
}