package cache

import (
	"time"
)

//go:generate moq  -out ../../mocks/mock_cache.go -skip-ensure -pkg mocks . Cache
type (
	// Config - for local cache.
	Config struct {
		DefaultTTL string `yaml:"default_ttl"`

		CustomCacheConfig map[string]any `yaml:"customCacheConfig"`
	}

	// TTL - struct describe TTL use cases for cache.
	TTL = struct {
		TTL      time.Duration // PREFER(more important) THAN ExpireAt
		ExpireAt time.Time
	}

	// Cache - interface of cache.
	Cache[K comparable, V any] interface {
		// Get - get data(Value) for Key (return Value if key doesn't expire)
		Get(k K) (V, bool)

		// Set - set data(Value) for Key (without TTL, not expired)
		Set(K, V)

		// SetWithTTL - set data(Value) for Key (with TTL or when expired)
		SetWithTTL(K, V, TTL)

		// SetTTL - set TTL for Key
		SetTTL(K, TTL)

		// TryGetOrInvokeLambda - try Get(Key) if not found, exec Lambda and if success store Value by Set func
		TryGetOrInvokeLambda(K, func(K) (V, TTL, error)) (V, error)

		// Delete - delete Key
		Delete(K)

		// CleanAll - delete all data in storage (all Key)
		CleanAll()
	}
)
