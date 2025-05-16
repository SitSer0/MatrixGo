package cache

import (
	"MatrixGo/internal/field"
	"MatrixGo/internal/matrix"
	"MatrixGo/internal/vector"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMatrixCache(t *testing.T) {
	// Создаем кэш с небольшим TTL для тестирования
	cache := NewMatrixCache[field.Float64](100*time.Millisecond, 10)

	// Создаем тестовую матрицу
	data := [][]field.Float64{
		{1, 2},
		{3, 4},
	}
	mat, _ := matrix.FromSlice(data)

	t.Run("cache determinant", func(t *testing.T) {
		// Первый вызов - вычисление
		start := time.Now()
		det1 := cache.CacheDeterminant(mat)
		firstCall := time.Since(start)

		// Второй вызов - должен быть из кэша
		start = time.Now()
		det2 := cache.CacheDeterminant(mat)
		secondCall := time.Since(start)

		// Проверяем результаты
		assert.Equal(t, det1, det2)
		assert.Less(t, secondCall, firstCall)

		// Ждем истечения TTL
		time.Sleep(150 * time.Millisecond)

		// Третий вызов - снова вычисление
		start = time.Now()
		det3 := cache.CacheDeterminant(mat)
		thirdCall := time.Since(start)

		assert.Equal(t, det1, det3)
		assert.Greater(t, thirdCall, secondCall)
	})

	t.Run("cache rank", func(t *testing.T) {
		// Первый вызов - вычисление
		start := time.Now()
		rank1 := cache.CacheRank(mat)
		firstCall := time.Since(start)

		// Второй вызов - должен быть из кэша
		start = time.Now()
		rank2 := cache.CacheRank(mat)
		secondCall := time.Since(start)

		// Проверяем результаты
		assert.Equal(t, rank1, rank2)
		assert.Less(t, secondCall, firstCall)
	})

	t.Run("cache inverse", func(t *testing.T) {
		// Первый вызов - вычисление
		start := time.Now()
		inv1, err1 := cache.CacheInverse(mat)
		firstCall := time.Since(start)
		assert.NoError(t, err1)

		// Второй вызов - должен быть из кэша
		start = time.Now()
		inv2, err2 := cache.CacheInverse(mat)
		secondCall := time.Since(start)
		assert.NoError(t, err2)

		// Проверяем результаты
		assert.Equal(t, inv1.Data, inv2.Data)
		assert.Less(t, secondCall, firstCall)
	})

	t.Run("cache solve system", func(t *testing.T) {
		vec := vector.NewVector([]field.Float64{1, 2})

		// Первый вызов - вычисление
		start := time.Now()
		sol1, err1 := cache.CacheSolveSystem(mat, vec)
		firstCall := time.Since(start)
		assert.NoError(t, err1)

		// Второй вызов - должен быть из кэша
		start = time.Now()
		sol2, err2 := cache.CacheSolveSystem(mat, vec)
		secondCall := time.Since(start)
		assert.NoError(t, err2)

		// Проверяем результаты
		assert.Equal(t, sol1.Data, sol2.Data)
		assert.Less(t, secondCall, firstCall)
	})
}

func TestCacheEviction(t *testing.T) {
	// Создаем кэш с максимальным размером 2 элемента
	cache := NewCache(1*time.Hour, 2)

	// Добавляем 3 элемента
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Проверяем, что один из старых элементов был удален
	count := 0
	if _, ok := cache.Get("key1"); ok {
		count++
	}
	if _, ok := cache.Get("key2"); ok {
		count++
	}
	if _, ok := cache.Get("key3"); ok {
		count++
	}

	assert.Equal(t, 2, count)
}

func TestCacheTTL(t *testing.T) {
	// Создаем кэш с TTL 100ms
	cache := NewCache(100*time.Millisecond, 10)

	// Добавляем значение
	cache.Set("key", "value")

	// Проверяем, что значение доступно
	val, ok := cache.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "value", val)

	// Ждем истечения TTL
	time.Sleep(150 * time.Millisecond)

	// Проверяем, что значение удалено
	_, ok = cache.Get("key")
	assert.False(t, ok)
}

func TestCacheThreadSafety(t *testing.T) {
	cache := NewCache(1*time.Hour, 1000)
	done := make(chan bool)

	// Запускаем несколько горутин для записи
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d-%d", id, j)
				cache.Set(key, j)
			}
			done <- true
		}(i)
	}

	// Запускаем несколько горутин для чтения
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d-%d", id, j)
				cache.Get(key)
			}
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 20; i++ {
		<-done
	}
}
