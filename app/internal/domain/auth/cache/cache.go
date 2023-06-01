package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) interface{}
	Set(key string, value string)
}

type cacheItem struct {
	Key       string
	Value     string
	AccessCnt int
}

type LRU struct {
	L1         map[string]*list.Element
	L1List     *list.List
	MaxSizeL1  int
	AccessTime time.Duration
	Lock       sync.Mutex
}

func NewCache(maxSizeL1 int, accessTime time.Duration) (Cache, error) {
	return &LRU{
		L1:         make(map[string]*list.Element),
		L1List:     list.New(),
		MaxSizeL1:  maxSizeL1,
		AccessTime: accessTime,
	}, nil
}

func (c *LRU) Get(key string) interface{} {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if item, ok := c.L1[key]; ok {
		c.promoteL1(item)
		return item.Value.(*cacheItem).Value
	}
	return nil
}

func (c *LRU) Set(key string, value string) {
	c.Lock.Lock()
	defer c.Lock.Unlock()
	if item, ok := c.L1[key]; ok {
		item.Value.(*cacheItem).Value = value
		c.promoteL1(item)
		return
	}

	item := &cacheItem{
		Key:       key,
		Value:     value,
		AccessCnt: 0,
	}

	if len(c.L1) < c.MaxSizeL1 {
		c.addToL1(item)
		return
	}

	c.evictL1()
	c.addToL1(item)
}

//func (c *LRU) Print() {
//	fmt.Println("Cache Contents:")
//	for key, item := range c.L1 {
//		fmt.Printf("Key: %s, Value: %v\n", key, item.Value.(*cacheItem).Value)
//	}
//}

func (c *LRU) promoteL1(item *list.Element) {
	item.Value.(*cacheItem).AccessCnt++
	item.Value.(*cacheItem).AccessCnt = 0
}

func (c *LRU) addToL1(item *cacheItem) {
	e := c.L1List.PushFront(item)
	c.L1[item.Key] = e
}

func (c *LRU) evictL1() {
	item := c.L1List.Back()
	c.L1List.Remove(item)
	delete(c.L1, item.Value.(*cacheItem).Key)
}
