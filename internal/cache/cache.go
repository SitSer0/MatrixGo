package cache

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// CacheEntry представляет запись в кэше
type CacheEntry struct {
	Value     interface{}
	Timestamp time.Time
}

// Cache представляет кэш для матричных операций
type Cache struct {
	data     map[string]CacheEntry
	mu       sync.RWMutex
	ttl      time.Duration
	maxItems int
}

// NewCache создает новый кэш
func NewCache(ttl time.Duration, maxItems int) *Cache {
	return &Cache{
		data:     make(map[string]CacheEntry),
		ttl:      ttl,
		maxItems: maxItems,
	}
}

// generateKey генерирует ключ для кэша на основе операции и входных данных
func generateKey(operation string, args ...interface{}) string {
	key := fmt.Sprintf("%s:%v", operation, args)
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// cleanup удаляет устаревшие записи
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, v := range c.data {
		if now.Sub(v.Timestamp) > c.ttl {
			delete(c.data, k)
		}
	}
}

// Set добавляет значение в кэш
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если кэш переполнен, удаляем самую старую запись
	if len(c.data) >= c.maxItems {
		var oldestKey string
		var oldestTime time.Time
		first := true

		for k, v := range c.data {
			if first || v.Timestamp.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.Timestamp
				first = false
			}
		}
		delete(c.data, oldestKey)
	}

	c.data[key] = CacheEntry{
		Value:     value,
		Timestamp: time.Now(),
	}
}

// Get получает значение из кэша
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Проверяем TTL
	if time.Now().Sub(entry.Timestamp) > c.ttl {
		go func() {
			c.mu.Lock()
			delete(c.data, key)
			c.mu.Unlock()
		}()
		return nil, false
	}

	return entry.Value, true
}

// MatrixCache представляет специализированный кэш для матричных операций
type MatrixCache[T field.Field[T]] struct {
	cache *Cache
}

// NewMatrixCache создает новый матричный кэш
func NewMatrixCache[T field.Field[T]](ttl time.Duration, maxItems int) *MatrixCache[T] {
	return &MatrixCache[T]{
		cache: NewCache(ttl, maxItems),
	}
}

// CacheDeterminant кэширует вычисление определителя
func (mc *MatrixCache[T]) CacheDeterminant(m *matrix.Matrix[T]) T {
	key := generateKey("det", m.Data)
	if val, ok := mc.cache.Get(key); ok {
		return val.(T)
	}

	result := m.DeterminantParallel()
	mc.cache.Set(key, result)
	return result
}

// CacheRank кэширует вычисление ранга
func (mc *MatrixCache[T]) CacheRank(m *matrix.Matrix[T]) int {
	key := generateKey("rank", m.Data)
	if val, ok := mc.cache.Get(key); ok {
		return val.(int)
	}

	result := m.RankParallel()
	mc.cache.Set(key, result)
	return result
}

// CacheInverse кэширует вычисление обратной матрицы
func (mc *MatrixCache[T]) CacheInverse(m *matrix.Matrix[T]) (*matrix.Matrix[T], error) {
	key := generateKey("inverse", m.Data)
	if val, ok := mc.cache.Get(key); ok {
		return val.(*matrix.Matrix[T]), nil
	}

	result, err := m.InverseParallel()
	if err != nil {
		return nil, err
	}

	mc.cache.Set(key, result)
	return result, nil
}

// CacheSolveSystem кэширует решение системы уравнений
func (mc *MatrixCache[T]) CacheSolveSystem(m *matrix.Matrix[T], v *vector.Vector[T]) (*vector.Vector[T], error) {
	key := generateKey("solve", m.Data, v.Data)
	if val, ok := mc.cache.Get(key); ok {
		return val.(*vector.Vector[T]), nil
	}

	result, err := matrix.SolveSystemParallel(m, v)
	if err != nil {
		return nil, err
	}

	mc.cache.Set(key, result)
	return result, nil
}
