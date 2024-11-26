package group

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/juroberttyb/group/models"
	"github.com/stretchr/testify/require"
)

// TODO: better test error handling, err encountered in go routines spawned in tests should cause fatal and stop testing from continuing
// TODO: make interface for group package and use mockery for testing

func TestConcurrentSameKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	resultCh := make(chan any)
	task := func(ctx context.Context, key string) (any, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- val
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- val
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
			runMsg,
			"result",
			"result",
		},
		result,
		"The two execution results should be the same.",
	)
}

func TestConcurrentDiffKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	resultCh := make(chan any)
	task := func(ctx context.Context, key string) (any, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- val
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "bar", task)
		t.Log(val, err)
		resultCh <- val
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
			runMsg,
			runMsg,
			"result",
			"result",
		},
		result,
		"The two execution results should be the same.",
	)
}

func TestSequentialSameKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: time.Second,
	})

	result := []any{}
	wg := sync.WaitGroup{}
	task := func(ctx context.Context, key string) (any, error) {
		defer wg.Done()
		t.Log(runMsg)
		result = append(result, runMsg)
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}

	wg.Add(1)
	val, err := group.Do(ctx, "foo", task)
	t.Log(val, err)
	result = append(result, val)

	wg.Add(1)
	val, err = group.Do(ctx, "foo", task)
	t.Log(val, err)
	result = append(result, val)

	wg.Wait()

	require.Equal(
		t,
		[]any{
			runMsg,
			"result",
			runMsg,
			"result",
		},
		result,
		"The two execution results should be the same.",
	)
}

// concurrently run `Group.Do()` with different keys with long running function
// output
// <nil> "timeout"
// <nil> "timeout"
// run
// run
func TestConcurrentLongRunDiffKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			Timeout: time.Second,
		},
	)

	resultCh := make(chan any)
	wg := sync.WaitGroup{}
	task := func(ctx context.Context, key string) (interface{}, error) {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		t.Log(runMsg)
		resultCh <- runMsg
		return "result", nil
	}

	wg.Add(1)
	go func() {
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- err.Error()
	}()

	wg.Add(1)
	go func() {
		val, err := group.Do(ctx, "bar", task)
		t.Log(val, err)
		resultCh <- err.Error()
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
			"timeout",
			"timeout",
			runMsg,
			runMsg,
		},
		result,
		"The two execution results should be the same.",
	)
}

func TestConcurrentTotalLimitDiffKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			MaxInflight: 2,
		},
	)

	resultCh := make(chan any)
	task := func(ctx context.Context, key string) (interface{}, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(time.Second)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log(val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "bar", task)
		t.Log(val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "baz", task)
		t.Log(val, err)
		resultCh <- err
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
		5,
		len(result),
		"Result length should be equal to 5 exactly.",
	)

	c := 0
	for _, e := range result {
		if e == models.ErrReachedLimit {
			c++
			break
		}
	}
	require.Equal(
		t,
		1,
		c,
		"Should have exactly least one max inflight limit error.",
	)
}

func TestConcurrentPerKeyLimitDiffKey(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](
		Options{
			MaxInflightPerKey: 2,
		},
	)
	resultCh := make(chan any)

	task := func(ctx context.Context, key string) (interface{}, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(time.Second)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log("foo 0", val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log("foo 1", val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task)
		t.Log("foo 2", val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "bar", task)
		t.Log("bar 0", val, err)
		resultCh <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "bar", task)
		t.Log("bar 1", val, err)
		resultCh <- err
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
		7,
		len(result),
		"Result length should be equal to 7 exactly.",
	)

	c := 0
	for _, e := range result {
		if e == models.ErrReachedLimitPerKey {
			c++
			break
		}
	}
	require.Equal(
		t,
		1,
		c,
		"Should have exactly one max inflight limit per key error.",
	)
}

func TestHashDontCollide(t *testing.T) {
	var runMsg = "run"

	ctx := context.Background()
	group := NewGroup[string, any](
		Options{},
	)
	resultCh := make(chan any)

	task0 := func(ctx context.Context, key string) (interface{}, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(time.Second)
		return "result", nil
	}

	task1 := func(ctx context.Context, key string) (interface{}, error) {
		t.Log(runMsg)
		resultCh <- runMsg
		time.Sleep(time.Second)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task0)
		t.Log("foo 0", val, err)
		resultCh <- val
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		val, err := group.Do(ctx, "foo", task1)
		t.Log("foo 1", val, err)
		resultCh <- val
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
			runMsg,
			runMsg,
			"result",
			"result",
		},
		result,
		"The two execution results should be the same.",
	)
}
