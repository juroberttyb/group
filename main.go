package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/juroberttyb/group/group"
)

func someTask(ctx context.Context, key string) (any, error) {
	time.Sleep(1 * time.Second)
	fmt.Println("run")
	return "result", nil
}

func main() {
	ctx := context.Background()
	wg := sync.WaitGroup{}
	wg.Add(2)
	group := group.NewGroup[string, any](group.Options{
		Timeout: 1 * time.Second,
	})
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", someTask))
	}()
	go func() {
		defer wg.Done()
		fmt.Println(group.Do(ctx, "foo", someTask))
	}()
	wg.Wait()
}

// output (only run someTask once)
// run
// result <nil>
// result <nil>
// -----------------------------------------------------------
// concurrently run `Group.Do()` with different keys
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))
// 	}()
// go func() {
//     fmt.Println(Group.Do(ctx, "bar", someTask))}()
// output
// runrunresult <nil>
// result <nil>
// -----------------------------------------------------------
// sequentially run `Group.Do()` with the same key
// fmt.Println(Group.Do(ctx, "foo", someTask))
// fmt.Println(Group.Do(ctx, "foo", someTask))
// output
// runresult <nil>
// runresult <nil>
// -----------------------------------------------------------
// concurrently run `Group.Do()` with different keys with long running function
// func someTask(ctx contex.Context, key string) (interface{}, error) {
//     time.Sleep(10 * time.Second)
//     return "result", nil
// }
// group := NewGroup(Options{Timeout: 5 * time.Second})
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "bar", someTask))}()
// output
// <nil> "timeout"
// <nil> "timeout"
// runrun
// -----------------------------------------------------------
// concurrently run `Group.Do()` with different keys and reached limit
// func someTask(ctx contex.Context, key string) (interface{}, error) {
//     time.Sleep(1 * time.Second)
//     fmt.Println("run")
//     return "result", nil
// }
// group := NewGroup(Options{MaxInflight: 2})
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "bar", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "baz", someTask))}()
// output
// runresult <nil>
// run
// result <nil>
// nil "reached inflight limit"
// -----------------------------------------------------------
// concurrently run `Group.Do()` with different keys and reached per key limit
// func someTask(ctx contex.Context, key string) (interface{}, error) {
//     time.Sleep(1 * time.Second)
//     fmt.Println("run", key)
//     return "result", nil
// }
// group := NewGroup(Options{MaxInflightPerKey: 2})
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "foo", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "bar", someTask))}()
// go func() {
//     fmt.Println(Group.Do(ctx, "bar", someTask))}()
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
