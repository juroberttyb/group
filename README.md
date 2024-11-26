# Group

This lib implements logic level caching by using func ```f()``` and its params ```x``` as key, the value is cached in ```f(x) <-> y``` mapping, where ```y=f(x)```.

Consecutive calls to ```f(x)``` receive result from the first call instead of running again.

## How to use

##### please run following command to test the lib
```
make unit-test
```

## Changes to testcase

### Add

Additional testcase ```TestHashDontCollide``` for checking whether ```f(x)``` hashes are generated correctly.

### Update

Some test cases are modified for better demo experience, for example

- Testcase: ```TestConcurrentSameKey```, ```TestConcurrentDiffKey```, and ```TestSequentialSameKey```
    
    Task is updated from

    ```
	task := ... {
		time.Sleep(time.Second)
		fmt.Println("run")
		return "result", nil
	}
    ``` 
    
    to

    ```
	task := ... {
		fmt.Println("run")
		time.Sleep(100 * time.Millisecond)
		return "result", nil
	}
    ```

    for better demo experience.

- Testcase: ```TestConcurrentTotalLimitDiffKey```
    
    Task is updated from

    ```
	task := ... {
		time.Sleep(time.Second)
		fmt.Println("run")
		return "result", nil
	}
    ``` 
    
    to

    ```
	task := ... {
		fmt.Println("run")
		time.Sleep(time.Second)
		return "result", nil
	}
    ```

    for better demo experience.

- Testcase: ```TestConcurrentTotalLimitDiffKey```

  Expected output is updated from

    ```
    run
    result <nil>
    run
    result <nil>
    <nil> reached inflight limit
    ```

    to

    ```
    run
    run
    <nil> reached inflight limit
    result <nil>
    result <nil>
    ```

    for better demo experience.

- Testcase: ```TestConcurrentPerKeyLimitDiffKey```

  Task is updated from

    ```
	task := ... {
        time.Sleep(time.Second)
		fmt.Println("run", key)
		return "result", nil
	}
    ```

  to 

    ```
	task := ... {
		fmt.Println("run", key)
		time.Sleep(time.Second)
		return "result", nil
	}
    ```

  and expected output is updated from

    ```
    run "foo"
    result <nil>
    run "foo"
    result <nil>
    nil "reached inflight limit"
    run "bar"
    result <nil>
    run "bar"
    result <nil>
    ```

    to

    ```
    run foo
    run bar
    foo 2 <nil> reached inflight limit
    bar 1 result<nil>
    bar 0 result<nil>
    foo 0 result<nil>
    foo 1 result<nil>
    ```

    for better demo experience.

- Testcase: ```TestConcurrentLongRunDiffKey```

    timeout is scaled down from

    ```
    5  seconds
    10 seconds
    ```

    to

    ```
    1  seconds
    2  seconds
    ```

    for faster execution result.