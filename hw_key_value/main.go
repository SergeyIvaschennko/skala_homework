package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	FirstReadTimeout  = 1 * time.Second
	SecondReadTimeout = 2 * time.Second
)

var CacheValues = [7]int{500, 200, 2000, 5000, 200, 1500, 5000}

type cacheItem struct {
	value     any
	expireAt  time.Time //19:55
	isEternal bool
}

type InMemoryCache struct {
	storage map[string]cacheItem
	mutex   sync.RWMutex

	cleanupTicker *time.Ticker //chan
	stopChannel   chan struct{}
}

func NewCache() *InMemoryCache {
	return &InMemoryCache{
		storage:     make(map[string]cacheItem),
		stopChannel: make(chan struct{}),
	}
}

func (cache *InMemoryCache) Run() {
	cache.cleanupTicker = time.NewTicker(500 * time.Millisecond) //1

	go func() {
		for {
			select {
			case <-cache.cleanupTicker.C: //1
				cache.cleanupExpired()
			case <-cache.stopChannel: //2
				cache.cleanupTicker.Stop()
				return
			}
		}
	}()
}

func (cache *InMemoryCache) Stop() { //2
	close(cache.stopChannel)
}

func (cache *InMemoryCache) Set(key string, value any, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	//immortal
	if ttl <= 0 {
		cache.storage[key] = cacheItem{
			value:     value,
			isEternal: true,
		}
		return
	}

	expirationTime := time.Now().Add(ttl)

	cache.storage[key] = cacheItem{
		value:    value,
		expireAt: expirationTime,
	}
}

func (cache *InMemoryCache) Get(key string) (any, bool) {
	cache.mutex.RLock()
	item, exists := cache.storage[key]
	cache.mutex.RUnlock()

	if !exists {
		return nil, false
	}

	if item.isEternal {
		return item.value, true
	}

	//check whether key is expired
	if time.Now().After(item.expireAt) {
		cache.Delete(key)
		return nil, false
	}

	return item.value, true
}

func (cache *InMemoryCache) Delete(key string) {
	cache.mutex.Lock()
	delete(cache.storage, key)
	cache.mutex.Unlock()
}

func (cache *InMemoryCache) cleanupExpired() {
	now := time.Now()

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, item := range cache.storage {
		if item.isEternal {
			continue
		}

		if now.After(item.expireAt) {
			delete(cache.storage, key)
		}
	}
}

func main() {
	cache := NewCache()
	cache.Run()
	defer cache.Stop()

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	// Запись значений
	go func() {
		defer waitGroup.Done()
		for index, value := range CacheValues {
			cache.Set(
				fmt.Sprintf("key_%d", index),
				value,
				time.Duration(value)*time.Millisecond,
			)
		}
	}()

	// Первое чтение
	go func() {
		defer waitGroup.Done()
		time.Sleep(FirstReadTimeout)

		for index := range CacheValues {
			key := fmt.Sprintf("key_%d", index)

			value, exists := cache.Get(key)
			if !exists {
				fmt.Printf("%s deleted\n", key)
				continue
			}

			fmt.Printf("%s: %v\n", key, value)
		}
	}()

	waitGroup.Wait()

	fmt.Println("\n      Pause       ")
	fmt.Println("\n      Pause       ")
	fmt.Println("\n      Pause       ")
	time.Sleep(SecondReadTimeout)

	for index := range CacheValues {
		key := fmt.Sprintf("key_%d", index)

		value, exists := cache.Get(key)
		if !exists {
			fmt.Printf("%s deleted\n", key)
			continue
		}

		fmt.Printf("%s: %v\n", key, value)
	}

	// Тест вечного ключа
	cache.Set("forever", "I live forever", 0)
	time.Sleep(3 * time.Second)

	if value, exists := cache.Get("forever"); exists {
		fmt.Printf("\nВечный ключ: %v\n", value)
	}

	// Тест удаления
	cache.Set("todelete", "delete me", 10*time.Second)
	cache.Delete("todelete")

	if _, exists := cache.Get("todelete"); !exists {
		fmt.Println("Ключ 'todelete' успешно удален")
	}

	// Тест конкурентности
	var concurrentWaitGroup sync.WaitGroup
	concurrentWaitGroup.Add(100)

	for index := 0; index < 100; index++ {
		currentIndex := index

		go func() {
			defer concurrentWaitGroup.Done()

			key := fmt.Sprintf("concurrent_%d", currentIndex)

			cache.Set(key, currentIndex, 5*time.Second)
			cache.Get(key)
		}()
	}

	concurrentWaitGroup.Wait()
	fmt.Println("Конкурентные операции завершены")
}
