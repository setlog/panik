# panik ![](https://github.com/setlog/panik/workflows/Tests/badge.svg)

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
    panik.OnErrorf(err, "f1() failed")
    b, err := f2(a)
    panik.OnErrorf(err, "f2(%v) failed", a)
    c, err := f3(b)
    panik.OnErrorf(err, "f3(%v) failed", b)
    return c
}

func getEverything() []interface{} {
    defer panik.Wrapf("could not get everything") // add more info to an ongoing panic
    s1 := getSomething()
    s2 := getSomethingElse()
    return []interface{} { s1, s2 }
}

func GetEverythingAndThenSome() (obj interface{}, retErr error) {
    defer panik.ToError(&retErr) // de-escalate panic into error
    return []interface{} { "and then some", getEverything()... }, nil
}

func iAmAGoroutine(everythingChannel chan interface{}) interface{} {
    defer panik.RecoverTrace() // if the panic could not be handled, end it all with some logging
    everythingChannel<-getEverything()
}

func iAmAnotherGoroutine() {
    defer func() {
        fmt.Println(recover() == nil) // false, because the panic did not originate from panik
    }()
    defer panik.Handle(func(r interface{}) { // will resume the panic at the end
        // clean up
    })
    panik.Wrap("error setting voltage level for flux compensator")
    panic("very critical problem")
}

func iAmYetAnotherGoroutine() {
    defer func() {
        fmt.Println(recover() == nil) // true, because the panic did originate from panik
    }()
    defer panik.Handle(func(r interface{}) { // recovers the panic
        // clean up
    })
    panik.Wrap("error processing item 42")
    panik.Panic("no biggie")
}

func getAnotherThing(id int) interface{} {
    if id == 42 {
        panik.Panicf("id %d is not supported", id) // panic from scratch when you have no non-nil error at hand
    }
    return things[id]
}
```

## Remarks
* Use `panik.ToError()` at API boundaries. APIs which panic are not idiomatic Go.
* You will still need to consider when to wrap an error and when to merely format its message using `%v`; the types of wrapped errors are part of your API. You can use `panik.IfError()` and `panik.IfErrorf()` for the highest level of control.
