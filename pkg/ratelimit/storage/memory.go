package storage

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type MemoryStorage struct {
	data        sync.Map // map[string]*memoryItem
	cleanUpMu   sync.Map
	stopCleanUp chan struct{}
}

type memoryItem struct {
	Value     interface{}
	ExpiredAt time.Time
}

// TTL cleanUp mechanism
// Define MemoryStorage.cleanUpLoop() method to execute background delete
func (ms *MemoryStorage) cleanUpLoop() {

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.cleanup()
		case <-ms.stopCleanUp:
			return
		}

	}
}

// Define MemoryStorage.cleanUp() method to delete expired data
func (ms *MemoryStorage) cleanup() {
	now := time.Now()
	keyToDelete := []string{}

	// check key to delete
	ms.data.Range(func(key, value interface{}) bool {
		mi := value.(*memoryItem)
		if !mi.ExpiredAt.IsZero() && now.After(mi.ExpiredAt) {
			keyToDelete = append(keyToDelete, key.(string))
		}
		return true
	})

	// delete key
	for _, key := range keyToDelete {
		ms.data.Delete(key)
	}

	if len(keyToDelete) > 0 {
		log.Printf("Cleaned up %d expired keys", len(keyToDelete))
	}

}

func (ms *MemoryStorage) Close() error {
	close(ms.stopCleanUp) // close the channel
	return nil
}

func NewMemoryStorage() *MemoryStorage {
	ms := &MemoryStorage{
		stopCleanUp: make(chan struct{}),
	}

	// clean up expired data in background whenever init
	go ms.cleanUpLoop()
	return ms
}

func (ms *MemoryStorage) Get(key string) (interface{}, error) {
	item, ok := ms.data.Load(key)
	if !ok {
		return nil, nil
	}

	mi := item.(*memoryItem)
	if !mi.ExpiredAt.IsZero() && time.Now().After(mi.ExpiredAt) {
		ms.data.Delete(key) // delete expired key
	}

	return mi.Value, nil
}

func (ms *MemoryStorage) Set(key string, value interface{}, ttl time.Duration) error {
	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = time.Now().Add(ttl)
	}
	ms.data.Store(key, &memoryItem{
		Value:     value,
		ExpiredAt: expiredAt,
	})

	return nil

}

func (ms *MemoryStorage) Increment(key string, ttl time.Duration) (int64, error) {
	var expiredAt time.Time
	if ttl > 0 {
		expiredAt = time.Now().Add(ttl)
	}
	for {
		item, loaded := ms.data.LoadOrStore(key, &memoryItem{
			Value:     int64(1),
			ExpiredAt: expiredAt,
		})
		mi := item.(*memoryItem)

		if !mi.ExpiredAt.IsZero() && time.Now().After(mi.ExpiredAt) {
			ms.data.Delete(key)
			continue
		}

		if loaded {
			currentValue, ok := mi.Value.(int64)
			if !ok {
				return 0, fmt.Errorf("value is not int64")
			}
			newValue := currentValue + 1
			ms.data.Store(key, &memoryItem{
				Value:     newValue,
				ExpiredAt: mi.ExpiredAt,
			})
			return newValue, nil
		}
		return 1, nil
	}
}

func (ms *MemoryStorage) Delete(key string) error {
	ms.data.Delete(key)
	return nil
}
