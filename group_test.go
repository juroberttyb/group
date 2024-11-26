package group

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// test whether only one function would run while there are two threads spawned with the same key
// expected output
// run
// result <nil>
// result <nil>
func TestConcurrentSameKey(t *testing.T) {

	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: 1 * time.Second,
	})

	task := func(ctx context.Context, key string) (any, error) {
		time.Sleep(1 * time.Second)
		fmt.Println("run")
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Wait()
}

// output
// run
// run
// result <nil>
// result <nil>
func TestConcurrentDiffKey(t *testing.T) {
	ctx := context.Background()
	group := NewGroup[string, any](Options{
		Timeout: 1 * time.Second,
	})

	task := func(ctx context.Context, key string) (any, error) {
		time.Sleep(1 * time.Second)
		fmt.Println("run")
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "bar", task))
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
		Timeout: 1 * time.Second,
	})

	task := func(ctx context.Context, key string) (any, error) {
		time.Sleep(1 * time.Second)
		fmt.Println("run")
		return "result", nil
	}

	fmt.Println(group.Do(ctx, "foo", task))
	fmt.Println(group.Do(ctx, "foo", task))
	time.Sleep(2 * time.Second)
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
		fmt.Println("run")
		return "result", nil
	}

	wg.Add(1)
	go func() {
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		fmt.Println(group.Do(ctx, "bar", task))
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
		time.Sleep(1 * time.Second)
		fmt.Println("run")
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "bar", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "baz", task))
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
		time.Sleep(1 * time.Second)
		fmt.Println("run", key)
		return "result", nil
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "bar", task))
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "bar", task))
	}()

	wg.Wait()
}
