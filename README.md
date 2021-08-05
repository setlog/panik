# panik ![](https://github.com/setlog/panik/workflows/Tests/badge.svg)

An error-handling paradigm for Go.

## The problem

```go
err := doSomething()
if err != nil {
    panic(err)
}
```

From the view of idiomatic Go, the above code has two problems:
1. An error is turned into a panic, making a treatable problem look like an untreatable one.
    * Even worse, it makes a treatable problem turn invisible, since a function signature does not inform about potential panics.
2. The error `err` will lack contextual information once it gets `recover()`ed, since no calls of the form `fmt.Errorf("-context-: %w", err)` are being made.

To alleviate these problems:
1. Call `recover()` before stepping over an API boundary.
    * If the recovered value isn't our error, assume the worst and `panic()` again.
2. There are errors where information on its circumstances aren't typically needed. Only use this technique with those. An example of this might be file IO: if the call fails, it is exceedingly likely that the cause is going to have to do with issues in the environment (disk full, missing permissions) which fall outside of the program's responsibility, or even control.

With those constraints in mind, we can help us out with a package such as this.

## panik's API and how to use it

The key-concept of panik is to reduce the amount of error-handling code by triggering a panic in places where you would otherwise return an error. As the stack unwinds, the panic may eventually be recovered (but only if it came from panik) to retrieve an error value again which then can be returned to the caller. This compares to Java's checked exceptions (i.e. an exception which hints at a problem in the environment as opposed to a problem within the program) except that there is no mechanism for declaring this in a function's signature in Go.

### panik by example

Most commonly, you start a pani*k* with `panik.Panic()` (if the function you are in decides that there is an error right now) or `panik.OnError(err)` (if a function-call returned an error to you). You then later end the pani*k* with `defer panik.ToError(&returnError)` (if you are at an API-boundary and don't want to give the caller some required reading) or one of the `defer panik.RecoverTrace...()` (if you are in some sort of worker) or `defer panik.ExitTrace...()` (if you are in `main()`) variants.

Panics triggered through panik are special in that panik can tell them apart from panics triggered through `panic()`, such that a call of `defer panik.ToError(&returnError)` will lastingly `recover()` and set `*returnError` only if the panic actually came from panik. This is accomplished with a package-private type which implements `error` as well as `Unwrap() error`. Basically, we make believe that a normal panic is a runtime exception while a "panik" is a checked exception.

```go
// doSomething(0) will return a nil error.
// doSomething(1) will return a non-nil error with Error() == "foo".
// doSomething(2) will panic.
func doSomething(x int) (returnError error) {
    defer panik.ToError(&returnError)
    doItNow(x)
    return nil
}

func doItNow(x int) {
    if x == 1 {
        panik.Panic("foo")
    } else if x == 2 {
        panic("foo")
    }
}
```

```go
// writeSomething("foo") will panic with a `ToError()`-deescalatable error if err is non-nil.
func writeSomething(filePath string) {
    err := ioutil.WriteFile(filePath, "Hello World!", 0660)
    panik.OnError(err)
}
```

### additional shenanigans

* You can use `defer panik.ToErrorWithTrace()` if your code is standalone and you really don't mind producing extra log output.
* You can use `defer panik.Wrap`(`f`)`()` to add extra information to an ongoing panic.
  * You will want to avoid calling this on a hot code path, since even if there is no panic, you are still making a function call with all of its arguments.
  * Using these functions violates point 2 as laid out in [The problem](#the-problem), so keep an eye out for whether panik is even appropriate for what you are doing.
* If calling `recover()` yourself, you can differ between panics and paniks using `panik.Caused(r)`.
  * You can always avoid having to do this by using `panik.ToError()` in the called function and then recovering in the caller normally.

## Remarks
* You will still need to think about when to wrap an error and when to merely format its message; the types of wrapped errors are part of your API contract. See `OnErrorfv()` if you have an error you want to report to the caller but do not want to wrap.
