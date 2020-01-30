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
    panik.IfErrorf(err, "f1() failed")
    b, err := f2(a)
    panik.IfErrorf(err, "f2(%v) failed", a)
    c, err := f3(b)
    panik.IfErrorf(err, "f3(%v) failed", b)
    return c
}

func getEverything() []interface{} {
    defer panik.Errorf("could not get everything: %w", panik.Cause{}) // add more info to an ongoing panic
    s1 := getSomething()
    s2 := getSomethingElse()
    return []interface{} { s1, s2 }
}

func GetEverythingAndThenSome() (obj interface{}, retErr error) {
    defer panik.ToError(&retErr) // de-escalate panic into error again
    return []interface{} { getEverything()..., "and then some" }, nil
}

func iAmAGoroutine(everythingChannel chan interface{}) interface{} {
    defer panik.RecoverTrace() // if the panic could not be handled, end it all with some logging
    everythingChannel<-getEverything()
}

// func iAmAnotherGoroutine() {
//     defer panik.WriteTrace(os.Stderr)
//     defer panik.Handle(func(r error) {
//         // never reached: plain panic() is not an error which is or wraps a *panik.knownCause.
//         // Only panik.Panicf() and panik.OnError() panic with such a value.
//     })
//     panic("very critical problem. DO NOT RECOVER")
// }

// func iAmYetAnotherGoroutine() { // a more explicit variant of iAmAnotherGoroutine
//     defer panik.WriteTrace(os.Stderr)
//     defer func() {
//         if r := recover(); r != nil {
//             if err, isError := r.(error); isError {
//                 var known *panik.Known
//                 if errors.As(err, &known) {
//                     fmt.Println("our code is aware of the origin of %v", r) // (a)
//                     return
//                 }
//             }
//             fmt.Println("our code has no idea where %v comes from", r) // (b)
//             panic(r)
//         }
//     }()

//     panic("(a)")
//     // OR
//     panik.Panic("(b)")
// }

func getAnotherThing(id int) interface{} {
    if id == 42 {
        panik.Panicf("id %d is not supported", id) // panic from scratch when you have no non-nil error at hand
    }
    return things[id]
}
```

## Remarks
* `err != nil` will continue to exist. You will also still want to perform type-assertions to get more information about an error's nature where appropriate. I.e. you cannot stop using your head here.
