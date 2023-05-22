package cache

import "container/list"

func (c *Cache) promoteL1(item *list.Element) {
	item.Value.(*CacheItem).AccessCnt++
	item.Value.(*CacheItem).AccessCnt = 0
}

func (c *Cache) addToL1(item *CacheItem) {
	e := c.L1List.PushFront(item)
	c.L1[item.Key] = e
}

func (c *Cache) evictL1() {
	item := c.L1List.Back()
	c.L1List.Remove(item)
	delete(c.L1, item.Value.(*CacheItem).Key)
}
