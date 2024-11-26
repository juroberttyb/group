package group

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
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
		inflight:     make(map[string]*call[T]),
		inflightKeys: make(map[string]int),
	}
}

// This return hash of arbitrary input to be used as inflight key
func hash(data []interface{}) (*string, error) {
	if len(data) == 0 {
		return nil, models.ErrListEmpty
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
	inflight      map[string]*call[T]
	inflightCount int
	inflightKeys  map[string]int
}

type call[T any] struct {
	wg  sync.WaitGroup
	val T
	err error
}

type Fx[K comparable, T any] struct {
	X K
	F Func[K, T]
}

type Func[K comparable, T any] func(ctx context.Context, key K) (value T, err error)

func (g *Group[K, T]) Do(ctx context.Context, key K, fn Func[K, T]) (value T, err error) {
	hkey, err := hash([]interface{}{
		key,
		reflect.ValueOf(fn).Pointer(),
	})
	if err != nil {
		return value, err
	}
	// println("hash info", hashKey, key, reflect.ValueOf(fn).Pointer())

	g.lock.Lock()
	// Check inflight limits
	if g.options.MaxInflight > 0 && g.inflightCount >= g.options.MaxInflight {
		g.lock.Unlock()
		return value, models.ErrReachedLimit
	}
	// Check inflight-key limits
	if g.options.MaxInflightPerKey > 0 && g.inflightKeys[*hkey] >= g.options.MaxInflightPerKey {
		g.lock.Unlock()
		return value, models.ErrReachedLimitPerKey
	}
	// Check if there's already an inflight call for this key
	if c, exists := g.inflight[*hkey]; exists {
		g.inflightCount++
		g.inflightKeys[*hkey]++
		g.lock.Unlock()

		// Wait for the in-flight call to finish
		c.wg.Wait()
		return c.val, c.err
	}

	// Register the new inflight call
	c := call[T]{}
	c.wg = sync.WaitGroup{}
	c.wg.Add(1)
	g.inflight[*hkey] = &c
	g.inflightCount++
	g.inflightKeys[*hkey]++
	g.lock.Unlock()

	// Execute the function
	defer func() {
		c.val = value
		c.err = err
		c.wg.Done()

		// Cleanup inflight tracking
		g.lock.Lock()
		defer g.lock.Unlock()
		delete(g.inflight, *hkey)
		g.inflightCount--
		g.inflightKeys[*hkey]--
	}()

	// Handle timeout
	if g.options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, g.options.Timeout)
		defer cancel()

		resCh := make(chan models.Result[T])
		go func() {
			val, err := fn(ctx, key)
			resCh <- models.Result[T]{
				Value: val,
				Err:   err,
			}
		}()

		// FIXME: is there a way to enforce fn to implement listener which receives kill signal?
		// a failed check in code review would cause infinite running threads to degrade server performance
		select {
		case <-ctx.Done():
			return value, models.ErrTimeout
		case res := <-resCh:
			return res.Value, res.Err
		}
	}

	value, err = fn(ctx, key)
	return
}
