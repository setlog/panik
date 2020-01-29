# panik

WORK IN PROGRESS (WIP)

An error-handling paradigm for Go.

The problem:

```go
func getSomething() (interface{}, error) {
    a, err := f1()
    if err != nil {
        return nil, fmt.Errorf("f1() failed: %w", err)
    }
    b, err := f2(a)
    if err != nil {
        return nil, fmt.Errorf("f2(%v) failed: %w", a, err)
    }
    c, err := f3(c)
    if err != nil {
        return nil, fmt.Errorf("f3(%v) failed: %w", b, err)
    }
    return c, nil
}
```

~75% of code is error handling. Also, there's a lot of it.

The solutions presented by this module, under the hood, make heavy use of Go 1.13 error-wrapping, `panic()` and `recover()`:

```go
func getSomething() interface{} { // panic instead of returning error
    a, err := f1()
    panik.OnError(err, "f1() failed") // ": %w" is automatically appended and filled using err when "%w" is missing.
    b, err := f2(a)
    panik.OnError(err, "f2(%v) failed", a) // fmt.Errorf()-style format args.
    c, err := f3(b)
    panik.OnError(err, "f3(%v) failed", b)
    return c
}
```

```go
func getEverything() []interface{} { // add more info to an ongoing panic
    defer panik.Described(err, "getSomething() failed") // topmost call must be Described() instead of Describe(). (See issues)
    s1 := getSomething()
    defer panik.Describe(err, "getSomethingElse() failed")
    s2 := getSomethingElse()
    return []interface{} { s1, s2 }
}
```

```go
func getEverythingAndThenSome() (obj interface{}, retErr error) { // de-escalate panic into error again
    defer panik.ToError(&retErr, "getEverything() failed") // on panic, set retErr to a non-nil error
    return []interface{} { getEverything()..., "and then some" }, nil
}
```

```go
func iAmAGoroutine(everythingChannel chan interface{}) interface{} {
    defer panik.WriteTrace(os.Stderr) // if the panic could not be handled, end it all with some logging
    defer panik.Handle(func(r error) { // handle known problems
        if r.Error() != "not too much of a problem" {
            panic(r)
        }
    })
    everythingChannel<-getEverything()
}

func iAmAnotherGoroutine() {
    defer panik.WriteTrace(os.Stderr)
    defer panik.Handle(func(r error) {
        // never reached: plain panic() is not an error which is or wraps a *panik.knownCause.
        // Only panik.Panic() and panik.OnError() panic with such a value.
    })
    panic("very critical problem. DO NOT RECOVER")
}

func iAmYetAnotherGoroutine() { // a more explicit variant of iAmAnotherGoroutine
    defer panik.WriteTrace(os.Stderr)
    defer func() {
        if r := recover(); r != nil {
            if err, isError := r.(error); isError {
                var known *panik.Known
                if errors.As(err, &known) {
                    fmt.Println("our code is aware of the origin of %v", r) // (a)
                    return
                }
            }
            fmt.Println("our code has no idea where %v comes from", r) // (b)
            panic(r)
        }
    }()

    panic("(a)")
    // OR
    panik.Panic("(b)")
}
```

```go
func getAnotherThing(id int) interface{} {
    if id == 42 {
        panik.Panic("id %d is not supported", id) // panic from scratch when you have no non-nil error at hand
    }
    return things[id]
}
```

## More

```go
func getSomething(somethingId int) (obj interface{}, retErr error) {
    defer panik.ToCustomError(&retErr, newIdError, somethingId) // de-escalate into your own implementation of the error interface
    return f(somethingId + 42), nil
}

func newIdError(cause error, args ...interface{}) error {
    return &IdError{cause: cause, id: args[0].(int)}
}

type IdError struct {
    cause error
    id    int
}

func (e *IdError) Error() string {
    return fmt.Sprintf("could not find id %d: %v", e.id, e.cause)
}

func (e *IdError) Unwrap() error {
    return e.cause
}

func (e *IdError) Id() int {
    return e.id
}
```

## Current Drawbacks
* `err != nil` will continue to exist. You will also still want to perform type-assertions to get more information about an error's nature where appropriate. I.e. you cannot stop using your head here.
* Forgetting to `panik.Described()` after `panik.Describe()` will leak memory.
