package group

import (
	"context"
	"fmt"
	"time"
)

func someTask(ctx context.Context, key string) (any, error) {
	fmt.Println("run")
	time.Sleep(1 * time.Second)
	return "result", nil
}

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
