# panik ![](https://github.com/setlog/panik/workflows/Tests/badge.svg)

An error-handling paradigm for Go.

## The problem

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

The solutions presented by this module, under the hood, make heavy use of Go 1.13 error-wrapping, `panic()` and `recover()`.

## panik's API and how to use it

The key-concept of panik is to reduce the amount of error-handling code by "converting" error values to panics. As the stack unwinds, the panic may eventually be recovered to retrieve an error value again which then can be returned to the caller.

### error to panic

#### when you have no error value at hand

```go
func Panic(args ...interface{}) {} // wraps errors.New(fmt.Sprint(args...))
func Panicf(format string, args ...interface{}) {} // wraps fmt.Errorf(format, args...)
```

#### when you have an error value or want to create one

```go
func OnError(err error) {} // this replaces if err != nil and panics with err when it is != nil
func OnErrorf(err error, format string, args ...interface{}) {} // panics with fmt.Errorf("%s: %w", fmt.Sprintf(format, args...)) with err appended to args when it is != nil
func IfError(err error, panicErr error) {} // calls OnError(panicErr) when err != nil
func IfErrorf(err error, format string, args ...interface{}) {} // shorthand for OnError(fmt.Errorf(format, args...)) when err != nil. In args, panik.Cause{} becomes err.
```

### panic to panic with more information

```go
func Wrap(args ...interface{}) {} // wraps value of panic in a new error with fmt.Sprint(args...) as its message if a panic is occurring.
func Wrapf(format string, args ...interface{}) {} // like Wrap, but with fmt.Sprintf(format, args...)
```

Use `Wrap()` and `Wrapf()` in `defer`-statements.

### panic to error

```go
func ToError(retErr *error) {} // assigns result of recover() to *retErr
```

Use `ToError()` in `defer`-statements.

### An important footnote

To prevent unwanted recovery and deescalation of panics originating from programming errors or from outside your own code, panik will never lastingly recover from panics created with `panic()`. Specifically, you need to use [one of panik's panicking functions](#error-to-panic) for `ToError()` to work.

## A practical overview

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
    defer panik.Recover(func(r interface{}) { // will resume the panic at the end
        // clean up
    })
    panik.Wrap("error setting voltage level for flux compensator")
    panic("very critical problem")
}

func iAmYetAnotherGoroutine() {
    defer func() {
        fmt.Println(recover() == nil) // true, because the panic did originate from panik
    }()
    defer panik.Recover(func(r interface{}) { // recovers the panic
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
