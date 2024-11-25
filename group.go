package group

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/juroberttyb/group/models"
)

type Options struct {
	// Timeout should applied no matter the input context is not done yet.
	// It'll return error `ErrTimeout` when tiemout.
	// `0` means no timeout.
	Timeout time.Duration
	// MaxInflight limit the number of running `Do` to the entire group.
	// It'll return `ErrReachedLimit` when the limit is reached.
	// `0` means no limit.
	MaxInflight int
	// MaxInflightPerKey limit the number of running `Do` to the entire group with same
	// It'll return `ErrReachedLimit` when the limit is reached.
	// `0` means no limit.
	MaxInflightPerKey int
}

func NewGroup[K comparable, T any](options Options) *Group[K, T] {
	return &Group[K, T]{
		options:      options,
		inflight:     make(map[K]*call[T]),
		inflightKeys: make(map[K]int),
	}
}

// TODO: verify data is hashable, think more on what might go wrong practically
// This return hash of arbitrary input to be used as inflight key
func MakeKey(data interface{}) (*string, error) {
	if data == nil {
		return nil, models.ErrNilInput
	}

	// Serialize the input data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input data: %w", err)
	}

	// Compute SHA-256 hash
	hash := sha256.Sum256(jsonData)
	hexHash := hex.EncodeToString(hash[:])

	// Convert the hash to a hex-encoded string
	return &hexHash, nil
}

type Group[K comparable, T any] struct {
	lock          sync.Mutex
	options       Options
	inflight      map[K]*call[T]
	inflightCount int
	inflightKeys  map[K]int
}

type call[T any] struct {
	wg  sync.WaitGroup
	val T
	err error
}

type Func[K comparable, T any] func(ctx context.Context, key K) (value T, err error)

func (g *Group[K, T]) Do(ctx context.Context, key K, fn Func[K, T]) (value T, err error) {
	println("function with key", key, "is in")

	g.lock.Lock()

	// Check inflight limits
	if g.options.MaxInflight > 0 && g.inflightCount >= g.options.MaxInflight {
		g.lock.Unlock()
		return value, models.ErrReachedLimit
	}
	// Check inflight-key limits
	if g.options.MaxInflightPerKey > 0 && g.inflightKeys[key] >= g.options.MaxInflightPerKey {
		g.lock.Unlock()
		return value, models.ErrReachedLimit
	}
	// Check if there's already an inflight call for this key
	if c, exists := g.inflight[key]; exists {
		g.inflightCount++
		g.inflightKeys[key]++
		g.lock.Unlock()
		println("function with key", key, "is waiting for peer to finish")

		// Wait for the in-flight call to finish
		c.wg.Wait()
		return c.val, c.err
	}

	// Register the new inflight call
	c := call[T]{}
	c.wg = sync.WaitGroup{}
	c.wg.Add(1)
	g.inflight[key] = &c
	g.inflightCount++
	g.inflightKeys[key]++
	g.lock.Unlock()

	// Execute the function
	defer func() {
		c.val = value
		c.err = err
		c.wg.Done()

		// Cleanup inflight tracking
		g.lock.Lock()
		defer g.lock.Unlock()
		delete(g.inflight, key)
		g.inflightCount--
		g.inflightKeys[key]--
	}()

	// Handle timeout
	if g.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.options.Timeout)
		defer cancel()
	}

	// FIXME: is there a way to enforce fn to implement
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 	default:
	// 		value, err = fn(ctx, key)
	// 	}
	// }
	// to listen to kill signal?
	// relying on only code review, a infinite running will degrade server performance
	value, err = fn(ctx, key)
	return
}