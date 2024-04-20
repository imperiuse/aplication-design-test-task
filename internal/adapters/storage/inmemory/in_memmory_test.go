package inmemory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	se "aplication-design-test-task/internal/adapters/storage"
)

func TestCreateDuplicateError(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage[int, string]()

	err := storage.Create(ctx, 1, "test")
	require.NoError(t, err)

	err = storage.Create(ctx, 1, "test")
	assert.ErrorIs(t, err, se.ErrDuplicateConstraint, "error should match ErrDuplicateConstraint when creating a duplicate item")
}

func TestReadNotFoundError(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage[int, string]()

	_, err := storage.Read(ctx, 1)
	assert.ErrorIs(t, err, se.ErrNotFound, "error should match ErrNotFound when reading a non-existent item")
}

func TestUpdateNotFoundError(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage[int, string]()

	err := storage.Update(ctx, 1, "update")
	assert.ErrorIs(t, err, se.ErrNotFound, "error should match ErrNotFound when updating a non-existent item")
}

func TestDeleteNotFoundError(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage[int, string]()

	err := storage.Delete(ctx, 1)
	assert.ErrorIs(t, err, se.ErrNotFound, "error should match ErrNotFound when deleting a non-existent item")
}

func TestCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()

	err := storage.Create(ctx, 1, "test")
	require.NoError(t, err, "should be able to create an item without error")

	err = storage.Create(ctx, 1, "test")
	assert.Error(t, err, "should not be able to create an item with duplicate ID")
}

func TestRead(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()

	require.NoError(t, storage.Create(ctx, 1, "test"))

	val, err := storage.Read(ctx, 1)
	require.NoError(t, err, "should be able to read an item without error")
	assert.Equal(t, "test", val, "the value should be 'test'")

	_, err = storage.Read(ctx, 2)
	assert.Error(t, err, "should not be able to read a non-existent item")
}

func TestUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()

	require.NoError(t, storage.Create(ctx, 1, "initial"))

	err := storage.Update(ctx, 1, "updated")
	require.NoError(t, err, "should be able to update an item without error")
	val, _ := storage.Read(ctx, 1)
	assert.Equal(t, "updated", val, "the value should be 'updated'")

	err = storage.Update(ctx, 2, "test")
	assert.Error(t, err, "should not be able to update a non-existent item")
}

func TestDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()

	require.NoError(t, storage.Create(ctx, 1, "test"))

	err := storage.Delete(ctx, 1)
	require.NoError(t, err, "should be able to delete an item without error")

	err = storage.Delete(ctx, 1)
	assert.Error(t, err, "should not be able to delete a non-existent item")
}

func TestList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()

	require.NoError(t, storage.Create(ctx, 1, "item1"))
	require.NoError(t, storage.Create(ctx, 2, "item2"))
	require.NoError(t, storage.Create(ctx, 3, "item3"))

	items, err := storage.List(ctx)
	require.NoError(t, err, "should be able to list items without error")
	assert.Len(t, items, 3, "should list all items correctly")
	expectedItems := []string{"item1", "item2", "item3"}
	assert.ElementsMatch(t, expectedItems, items, "the returned list should contain all added items")
}

func TestConcurrency(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage := NewInMemoryStorage[int, string]()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := storage.Create(ctx, id, "value")
			assert.NoError(t, err, "should be able to create items concurrently without error")
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			val, err := storage.Read(ctx, id)
			assert.NoError(t, err, "should be able to read items concurrently without error")
			assert.Equal(t, "value", val, "the value should be 'value'")
		}(i)
	}
	wg.Wait()

	// Concurrent updates
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := storage.Update(ctx, id, "updated")
			assert.NoError(t, err, "should be able to update items concurrently without error")
		}(i)
	}
	wg.Wait()

	// Concurrent deletes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := storage.Delete(ctx, id)
			assert.NoError(t, err, "should be able to delete items concurrently without error")
		}(i)
	}
	wg.Wait()
}

func TestCreateWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Ensure the context is cancelled to free resources if the test finishes before the timeout

	storage := NewInMemoryStorage[int, string]()

	// This goroutine simulates a delayed cancellation to affect the ongoing operation
	go func() {
		time.Sleep(500 * time.Millisecond) // Simulate some operation delay
		cancel()                           // Trigger cancellation
	}()

	// Wait to ensure the context has been cancelled
	<-ctx.Done()

	// Attempt to create an item after the context has been cancelled
	err := storage.Create(ctx, 1, "test")

	// Verify that the operation responds correctly to the context being cancelled
	assert.Error(t, err, "should return an error when the context is cancelled")
	assert.Equal(t, context.Canceled, err, "error should be context.Canceled")

	err = storage.Update(ctx, 1, "updated")
	assert.Error(t, err, "should return an error when the context is cancelled")
	assert.Equal(t, context.Canceled, err, "error should be context.Canceled")

	err = storage.Delete(ctx, 1)
	assert.Error(t, err, "should return an error when the context is cancelled")
	assert.Equal(t, context.Canceled, err, "error should be context.Canceled")

	_, err = storage.List(ctx)
	assert.Error(t, err, "should return an error when the context is cancelled")
	assert.Equal(t, context.Canceled, err, "error should be context.Canceled")

	_, err = storage.Read(ctx, 1)
	assert.Error(t, err, "should return an error when the context is cancelled")
	assert.Equal(t, context.Canceled, err, "error should be context.Canceled")
}
