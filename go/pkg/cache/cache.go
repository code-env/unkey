package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/maypok86/otter"
	"github.com/panjf2000/ants"
	"github.com/unkeyed/unkey/go/pkg/clock"
	"github.com/unkeyed/unkey/go/pkg/fault"
	"github.com/unkeyed/unkey/go/pkg/logging"
)

type cache[K comparable, V any] struct {
	otter otter.Cache[K, swrEntry[V]]
	fresh time.Duration
	stale time.Duration
	// If a key is stale, its key will be put into this channel and a goroutine refreshes it in the background
	refreshC chan K
	logger   logging.Logger
	resource string
	clock    clock.Clock

	inflightMu        sync.Mutex
	inflightRefreshes map[K]bool

	pool *ants.Pool
}

type Config[K comparable, V any] struct {
	// How long the data is considered fresh
	// Subsequent requests in this time will try to use the cache
	Fresh time.Duration

	// Subsequent requests that are not fresh but within the stale time will return cached data but also trigger
	// fetching from the origin server
	Stale time.Duration

	Logger logging.Logger

	// Start evicting the least recently used entry when the cache grows to MaxSize
	MaxSize int

	Resource string

	Clock clock.Clock
}

var _ Cache[any, any] = (*cache[any, any])(nil)

// New creates a new cache instance
func New[K comparable, V any](config Config[K, V]) *cache[K, V] {

	builder, err := otter.NewBuilder[K, swrEntry[V]](config.MaxSize)
	if err != nil {
		panic(err)
	}

	otter, err := builder.CollectStats().Cost(func(key K, value swrEntry[V]) uint32 {
		return 1
	}).WithTTL(time.Hour).Build()
	if err != nil {
		panic(err)
	}

	pool, err := ants.NewPool(10)
	if err != nil {
		panic(err)
	}

	c := &cache[K, V]{
		otter:             otter,
		fresh:             config.Fresh,
		stale:             config.Stale,
		refreshC:          make(chan K, 1000),
		logger:            config.Logger,
		resource:          config.Resource,
		clock:             config.Clock,
		pool:              pool,
		inflightMu:        sync.Mutex{},
		inflightRefreshes: make(map[K]bool),
	}

	return c

}

func (c *cache[K, V]) Get(ctx context.Context, key K) (value V, hit CacheHit) {

	e, ok := c.otter.Get(key)
	if !ok {
		// This hack is necessary because you can not return nil as V
		var v V
		return v, Miss
	}

	now := c.clock.Now()

	if now.Before(e.Stale) {
		return e.Value, e.Hit
	}

	c.otter.Delete(key)

	var v V
	return v, Miss

}

func (c *cache[K, V]) SetNull(ctx context.Context, key K) {
	c.set(ctx, key)
}

func (c *cache[K, V]) Set(ctx context.Context, key K, value V) {
	c.set(ctx, key, value)
}
func (c *cache[K, V]) set(_ context.Context, key K, value ...V) {
	now := c.clock.Now()

	e := swrEntry[V]{
		Value: value[0],
		Fresh: now.Add(c.fresh),
		Stale: now.Add(c.stale),
		Hit:   Null,
	}
	if len(value) > 0 {
		e.Value = value[0]
		e.Hit = Hit
	} else {
		e.Hit = Miss
	}
	c.otter.Set(key, e)

}

func (c *cache[K, V]) Remove(ctx context.Context, key K) {

	c.otter.Delete(key)

}

func (c *cache[K, V]) Dump(ctx context.Context) ([]byte, error) {
	data := make(map[K]swrEntry[V])

	c.otter.Range(func(key K, entry swrEntry[V]) bool {
		data[key] = entry
		return true
	})

	b, err := json.Marshal(data)

	if err != nil {
		return nil, fault.Wrap(err, fault.WithDesc("failed to marshal cache data", ""))
	}
	return b, nil

}

func (c *cache[K, V]) Restore(ctx context.Context, b []byte) error {

	data := make(map[K]swrEntry[V])
	err := json.Unmarshal(b, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}
	now := c.clock.Now()
	for key, entry := range data {
		if now.Before(entry.Fresh) {
			c.Set(ctx, key, entry.Value)
		} else if now.Before(entry.Stale) {
			c.refreshC <- key
		}
		// If the entry is older than, we don't restore it
	}
	return nil
}

func (c *cache[K, V]) Clear(ctx context.Context) {
	c.otter.Clear()
}

func (c *cache[K, V]) refresh(
	ctx context.Context,
	key K, refreshFromOrigin func(context.Context) (V, error),
	translateError func(error) CacheHit,
) {
	c.inflightMu.Lock()
	_, ok := c.inflightRefreshes[key]
	if ok {
		c.inflightMu.Unlock()
		return
	}
	c.inflightRefreshes[key] = true
	c.inflightMu.Unlock()

	defer func() {
		c.inflightMu.Lock()
		delete(c.inflightRefreshes, key)
		c.inflightMu.Unlock()
	}()

	v, err := refreshFromOrigin(ctx)

	switch translateError(err) {
	case Hit:
		c.set(ctx, key, v)
	case Miss:
		c.set(ctx, key)
	case Null:
		c.set(ctx, key)
	}

}

func (c *cache[K, V]) SWR(
	ctx context.Context,
	key K,
	refreshFromOrigin func(context.Context) (V, error),
	translateError func(error) CacheHit,
) (V, error) {
	now := c.clock.Now()
	e, ok := c.otter.Get(key)
	if ok {
		// Cache Hit

		if now.Before(e.Fresh) {
			// We have data and it's fresh, so we return it

			return e.Value, nil
		}

		if now.Before(e.Stale) {
			// We have data, but it's stale, so we refresh it in the background
			// but return the current value

			err := c.pool.Submit(func() {
				c.refresh(ctx, key, refreshFromOrigin, translateError)
			})
			if err != nil {
				c.logger.Error(ctx, "failed to submit refresh task", slog.String("error", err.Error()))
			}

			return e.Value, nil
		}

		// We have old data, that we should not serve anymore
		c.otter.Delete(key)

	}
	// Cache Miss

	// We have no data and need to go to the origin

	v, err := refreshFromOrigin(ctx)

	switch translateError(err) {
	case Hit:
		c.set(ctx, key, v)
	case Miss:
		c.set(ctx, key)
	case Null:
		c.set(ctx, key)
	}

	return v, err

}
