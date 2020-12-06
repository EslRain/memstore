package memory_store

import (
	"fmt"
	"sync"
	"time"
)

type MemStore struct {
	Store *sync.Map
}

type itemWithTTL struct {
	expires int64
	value interface{}
}

func NewMemStore() *MemStore {
	return &MemStore{
		Store: &sync.Map{},
	}
}

func newItem(value interface{}, expires int) itemWithTTL {
	expires64 := int64(expires)
	if expires64 > 0 {
		expires64 = time.Now().Unix() + expires64
	}

	return itemWithTTL{
		value: value,
		expires: expires64,
	}
}

func getValue(item interface{}, ok bool) (interface{}, bool) {
	if !ok {
		return nil, ok
	}

	var itemObj itemWithTTL
	if itemObj, ok = item.(itemWithTTL); !ok {
		return item, true
	}

	if itemObj.expires > 0 && itemObj.expires < time.Now().Unix() {
		return nil, false
	}

	return itemObj.value, ok
}

func (store *MemStore) Set(key string, value interface{}, ttl int) error {
	store.Store.Store(key, newItem(value, ttl))
	return nil
}

func (store *MemStore) Get(key string) (interface{}, bool) {
	//将垃圾回收操作平摊到每次Get操作中
	store.GarbageCollect()

	return getValue(store.Store.Load(key))
}

func (store *MemStore) GarbageCollect() {
	store.Store.Range(func(key, value interface{}) bool {
		if item, ok := value.(itemWithTTL); ok {
			if item.expires > 0 && item.expires < time.Now().Unix() {
				fmt.Printf("垃圾回收[%s]\n", key.(string))
				store.Store.Delete(key)
			}
		}
	})
}

