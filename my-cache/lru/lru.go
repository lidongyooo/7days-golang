package lru

import "container/list"

type Cache struct {
	cache map[string]*list.Element
	list *list.List
	maxBytes int64
	currBytes int64
	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int64
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		cache: make(map[string]*list.Element),
		list: list.New(),
		maxBytes: maxBytes,
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.list.MoveToFront(ele)
		return kv.value, true
	}

	return nil, false
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.currBytes -= kv.value.Len() - value.Len()
		kv.value = value

		c.list.MoveToFront(ele)
	} else {
		ele := c.list.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.currBytes += int64(len(key)) + value.Len()
	}

	for c.maxBytes != 0 && c.currBytes > c.maxBytes {
		c.Remove()
	}
}

func (c *Cache) Len() int {
	return c.list.Len()
}

func (c *Cache) Remove() {
	ele := c.list.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		c.currBytes -= int64(len(kv.key)) + kv.value.Len()
		c.list.Remove(ele)
		delete(c.cache, kv.key)

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
	
