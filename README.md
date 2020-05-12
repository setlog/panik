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

The solutions presented by this module make heavy use of Go 1.13 error-wrapping, `panic()` and `recover()` under the hood.

## panik's API and how to use it

The key-concept of panik is to reduce the amount of error-handling code by "converting" error values to panics. As the stack unwinds, the panic may eventually be recovered to retrieve an error value again which then can be returned to the caller.

### error to panic

#### when you have no error value at hand

```go
func Panic(r interface{}) {} // wraps panic()
func Panicf(format string, args ...interface{}) {} // wraps panic(fmt.Errorf(format, args...))
```

#### when you have an error at hand

The following functions only do something `if err != nil`, and act as a replacement for such.

```go
func OnError(err error) {} // panics with an error which wraps err
func OnErrore(err error, panicErr error) {} // panics with an error which wraps panicErr and returns fmt.Sprintf("%v: %v", panicErr, err) for Error()
func OnErrorfw(err error, format string, args ...interface{}) {} // panics with an error which wraps err and returns fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err) for Error()
func OnErrorfv(err error, format string, args ...interface{}) {} // panics with an error which does not wrap err and returns fmt.Sprintf("%s: %v", fmt.Sprintf(format, args...), err) for Error()
```

### panic to panic with more information

```go
func Wrap(args ...interface{}) {} // wraps value of panic in a new error with fmt.Sprint(args...) as its message if a panic is occurring.
func Wrapf(format string, args ...interface{}) {} // like Wrap, but with fmt.Sprintf(format, args...)
```

Use `Wrap()` and `Wrapf()` with the `defer`-statement.

### panic to error

```go
func ToError(retErr *error) {} // assigns result of recover() to *retErr
func ToErrorWithTrace(retErr *error) {} // like ToError(), but also contains full stack trace in the error's message
```

Use `ToError()` and `ToErrorWithTrace()` with the `defer`-statement.

#### an important footnote

To prevent unwanted recovery and deescalation of panics originating from programming errors, panik will never lastingly recover from panics created with `panic()`. Specifically, you need to use [one of panik's panicking functions](#error-to-panic) for `defer ToError(&err)` to set `*err`.

### inspect recovered panic value

```go
func Caused(r interface{}) bool {} // returns true if r was recovered from a panic started with panik.
```

### print a stack trace

```go
func RecoverTraceTo(w io.Writer) {} // to w
func RecoverTraceFunc(f func(trace string)) {} // to whatever else, e.g. an error dialog box (convenience function)
func ExitTraceTo(w io.Writer) {} // like RecoverTraceTo(), followed by os.Exit(2)
func ExitTraceFunc(f func(trace string)) {} // like RecoverTraceFunc(), followed by os.Exit(2)
```

Use `RecoverTrace`(`Func`)`()` in libraries. Use `ExitTrace`(`Func`)`()` in `main()`.

## A practical overview

```go
func getSomething() interface{} { // panic instead of returning error
    a, err := f1()
    panik.OnErrorfw(err, "f1() failed")
    b, err := f2(a)
    panik.OnErrorfw(err, "f2(%v) failed", a)
    c, err := f3(b)
    panik.OnErrorfw(err, "f3(%v) failed", b)
    return c
}

func getEverything() []interface{} {
    defer panik.Wrap("could not get everything") // add more info to an ongoing panic
    s1 := getSomething()
    s2 := getSomethingElse()
    return []interface{} { s1, s2 }
}

func GetEverythingAndThenSome() (obj interface{}, retErr error) {
	defer panik.ToError(&retErr) // de-escalate panic into error
	return append(getEverything(), "and then some"), nil
}

func iAmAGoroutine(everythingChannel chan interface{}) interface{} {
    defer panik.RecoverTraceTo(os.Stderr) // if the panic could not be handled, end it all with some logging
    everythingChannel<-getEverything()
}

func getAnotherThing(id int) interface{} {
    if id == 42 {
        panik.Panicf("cannot handle id %d", id) // panic from scratch when you have no non-nil error at hand
    }
    return things[id]
}

func DoSomething() (retErr error) {
    defer panik.ToError(&retErr)
    panik.OnErrorfv(doSomethingInternal(), "could not do something") // eliminate type information about type "superSpecificError".

}

func doSomethingInternal() error {
    // ...
    return &superSpecificError{}
}
```

## Remarks
* Use `panik.ToError()` at API boundaries. APIs which panic are not idiomatic Go.
    * Use `panik.ToErrorWithTrace()` if your error message is not informative enough by itself, but use it sparingly: stack traces are associated with programming errors; this will look bad if you don't know what the caller is going to do with the returned error.
* Avoid calling `recover()` yourself. If you do, you take the responsibility of differing between panics caused by your code and panics of unknown origin using `panik.Caused()`. Also, [panik's trace-printing functions](#print-a-stack-trace) will have no visibility of the original panic. You basically lose all of the simplicity panik provides. It is adviced to use `panik.ToError()` instead.
* You will still need to think about when to wrap an error and when to merely format its message; the types of wrapped errors are part of your API contract. See `OnErrorfv()` if you do not want to wrap an error.
