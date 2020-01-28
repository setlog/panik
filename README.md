# panik

WIP

An error-handling paradigm for Go.

The problem:

```go
func getSomething() (interface{}, error) {
    a, err := f1()
    if err != nil {
        return nil, err
    }
    b, err := f2(a)
    if err != nil {
        return nil, err
    }
    c, err := f3(c)
    if err != nil {
        return nil, err
    }
    return c, nil
}
```

~75% of code is error handling. Also, there's a lot of it.

The solution presented by this module:

```go
func getSomething() (interface{}, retErr error) {
    defer panik.ToError(&retErr, "could not get something")
    return f3(f2(f1)), nil // Here, f3, f2 and f1 panic on error
}
```

## More

```go
func getSomething() interface{} { // instead of f3(f2(f1))
    defer panik.Described("f1() failed") // Described() reenables describing for callers.
    a := f1()
    defer panik.Describe("f2(%v) failed", a)
    b := f2(a)
    defer panik.Describe("f3(%v) failed", b)
    return f3(b)
}
```

```go
func iAmAGoroutine(somethingChannel chan interface{}) interface{} {
    defer panik.WriteTrace(os.Stderr)
    somethingChannel<-getSomething()
}
```

```go
func iAmAGoroutine(somethingChannel chan interface{}) interface{} {
    defer panik.WriteTrace(os.Stderr)
    somethingChannel<-getSomething()
}
```

## Drawbacks
* By reinventing exceptions using `panic()` under the hood, one can end up recovering from more serious problems which one should not lastingly recover from. We recommend you only use this module with the intention of getting rid of panicking contexts.
* `err != nil` will continue to exist. Namely your lower-level functions/methods will have to check for it on whatever they are calling to react with a `panic()` instead of returning `err`. However, the total amount of error-checking will go down.
* Forgetting to `panik.Described()` after `panik.Describe()` will leak memory.
