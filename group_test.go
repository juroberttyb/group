package group

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TODO: make interface for group package and use mockery for testing
// TODO: all these tests require assert conditions

func TestConcurrentSameKey(t *testing.T) {
	ctx := context.Background()

	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	resultCh := make(chan any)
	task := func(ctx context.Context, key string) (any, error) {
		t.Log("run")
		resultCh <- "run"
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- val.(string)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- val.(string)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result := []any{}
	for {
		res, exist := <-resultCh
		if !exist {
			break
		}
		result = append(result, res)
	}

	require.Equal(
		t,
		[]any{
			"run",
			"result",
			"result",
		},
		result,
		"The two execution results should be the same.",
	)
}

// output
// run
// run
// result <nil>
// result <nil>
func TestConcurrentDiffKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	task := func(ctx context.Context, key string) (any, error) {
		t.Log("run")
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log(group.Do(ctx, "bar", task))
	}()

	wg.Wait()
}

// output
// run
// run
// result <nil>
// result <nil>
func TestSequentialSameKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	wg := sync.WaitGroup{}
	task := func(ctx context.Context, key string) (any, error) {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		t.Log("run")
		return "result", nil
	}

	wg.Add(1)
	t.Log(group.Do(ctx, "foo", task))

	wg.Add(1)
	t.Log(group.Do(ctx, "foo", task))

	wg.Wait()
}

// concurrently run `Group.Do()` with different keys with long running function
// output
// <nil> "timeout"
// <nil> "timeout"
// run
// run
func TestConcurrentLongRunDiffKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			Timeout: 5 * time.Second,
		},
	)

	wg := sync.WaitGroup{}
	task := func(ctx context.Context, key string) (interface{}, error) {
		defer wg.Done()
		time.Sleep(10 * time.Second)
		t.Log("run")
		return "result", nil
	}

	wg.Add(1)
	go func() {
		t.Log(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		t.Log(group.Do(ctx, "bar", task))
	}()

	wg.Wait()
}

// concurrently run `Group.Do()` with different keys and reached limit
// output
// run
// result <nil>
// run
// result <nil>
// nil "reached inflight limit"
func TestConcurrentTotalLimitDiffKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			MaxInflight: 2,
		},
	)

	task := func(ctx context.Context, key string) (interface{}, error) {
		t.Log("run")
		time.Sleep(time.Second)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log(group.Do(ctx, "bar", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log(group.Do(ctx, "baz", task))
	}()

	wg.Wait()
}

// concurrently run `Group.Do()` with different keys and reached per key limit
// output
// run "foo"
// result <nil>
// run "foo"
// result <nil>
// nil "reached inflight limit"
// run "bar"
// result <nil>
// run "bar"
// result <nil>
func TestConcurrentPerKeyLimitDiffKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			MaxInflightPerKey: 2,
		},
	)

	task := func(ctx context.Context, key string) (interface{}, error) {
		t.Log("run", key)
		time.Sleep(time.Second)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("foo 0",
			fmt.Sprint(group.Do(ctx, "foo", task)),
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("foo 1",
			fmt.Sprint(group.Do(ctx, "foo", task)),
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("foo 2",
			fmt.Sprint(group.Do(ctx, "foo", task)),
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("bar 0",
			fmt.Sprint(group.Do(ctx, "bar", task)),
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.Log("bar 1",
			fmt.Sprint(group.Do(ctx, "bar", task)),
		)
	}()

	wg.Wait()
}
