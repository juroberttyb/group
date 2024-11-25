Single Flight
Please implement func NewGroup, type Group and func (*Group) Do and corresponding unit tests.
Group.Do() makes sure that only one execution of the given function is in-flight for a given key. 
If another function call with the duplicated key comes in, the duplicated caller waits for the in-flight call with the same key to complete and receives the same results.

1.
Please implement func NewGroup, type Group and func (*Group) Do WITHOUT using any package except
Go builtin packages listed here https://pkg.go.dev/std
2.
Unit testing is a plus

var (
    ErrTimeout = errors.New("timeout")
    ErrReachedLimit = errors.New("reached inflight limit")
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

func NewGroup(options Options) *Group {
    // Please implement
}

type Group[K comparable, T any] struct {
    // Please implement
}

type Func[K comparable, T any] func(ctx context.Context, key K) (value T, err error)

func (s *Group[K, T]) Do(ctx context.Context, key K, fn Func[K, T]) (value T, err erro
    // Please implement
}

Group.Do() usage example:
package main

var (
    "fmt"
    "sync"
    "time"
)
func someTask(ctx contex.Context, key string) (any, error) {
    time.Sleep(1 * time.Second)
    fmt.Println("run")
    return "result", nil
}

// concurrently run `Group.Do()` with the same key
func main() {
    ctx := context.Background()
    var wg = sync.WaitGroup
    wg.Add(2)
    group := NewGroup[string, any](Options{
    Timeout: 1 * time.Second,
})

go func() {
    defer wg.Done()
    fmt.Println(group.Do(ctx, "foo", someTask))
}()

go func() {
    defet wg.Done()
    fmt.Println(group.Do(ctx, "foo", someTask))
}()

wg.Wait()
}
// output (only run someTask once)
// run
// result <nil>
// result <nil>
-----------------------------------------------------------
// concurrently run `Group.Do()` with different keys
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "bar", someTask))
}()
// output
run
run
result <nil>
result <nil>
-----------------------------------------------------------
// sequentially run `Group.Do()` with the same key
    fmt.Println(Group.Do(ctx, "foo", someTask))
    fmt.Println(Group.Do(ctx, "foo", someTask))
// output
run
result <nil>
run
result <nil>
-----------------------------------------------------------
// concurrently run `Group.Do()` with different keys with long running function
func someTask(ctx contex.Context, key string) (interface{}, error) {
    time.Sleep(10 * time.Second)
    return "result", nil
}
group := NewGroup(Options{Timeout: 5 * time.Second})
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "bar", someTask))
}()
// output
<nil> "timeout"
<nil> "timeout"
run
run
-----------------------------------------------------------
// concurrently run `Group.Do()` with different keys and reached limit
func someTask(ctx contex.Context, key string) (interface{}, error) {
    time.Sleep(1 * time.Second)
    fmt.Println("run")
    return "result", nil
}
group := NewGroup(Options{MaxInflight: 2})
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "bar", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "baz", someTask))
}()
// output
run
result <nil>
run
result <nil>
nil "reached inflight limit"
-----------------------------------------------------------
// concurrently run `Group.Do()` with different keys and reached per key limit
func someTask(ctx contex.Context, key string) (interface{}, error) {
    time.Sleep(1 * time.Second)
    fmt.Println("run", key)
    return "result", nil
}
group := NewGroup(Options{MaxInflightPerKey: 2})
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "foo", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "bar", someTask))
}()
go func() {
    fmt.Println(Group.Do(ctx, "bar", someTask))
}()
// output
run "foo"
result <nil>
run "foo"
result <nil>
nil "reached inflight limit"
run "bar"
result <nil>
run "bar"
result <nil>