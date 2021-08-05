# v1.0.0 (2021-08-05)
- Calling `panic()` with a panic-value from panik recovered with `defer ToError()` or `defer ToErrorWithTrace()` will no longer be recovered from by further calls of `defer ToError()` and `defer ToErrorWithTrace()`.
- `defer ToError()` now guarantees to extract the same error which was passed to `OnError()`.
- Rewrote README.md.
- Reversed order of entries in this changelog to be latest first.

# v0.5.1 (2021-07-02)
- Added `RecoverTraceToDefaultLogger()` and `ExitTraceToDefaultLogger()`.

# v0.5.0 (2020-05-12)
- Removed `RecoverTrace()` and `ExitTrace()`.
- Added `ToErrorWithTrace()`.

# v0.4.1 (2020-04-21)
- Added `RecoverTraceFunc()` and `ExitTraceFunc()`.

# v0.4.0 (2020-02-07)
- Added more info to README.md.
- Removed `Handle()`.
- `panik.value` is now `panik.Value`, allowing users to inspect the value.
- Changed signature of `Panic()` to be consistent with `panic()`.
- Simplified API.

# v0.3.0 (2020-01-31)
- Major rewrite.

# v0.2.0 (2020-01-29)
- Added `OnError()`.
- Added `Panic()`.
- Added `IsKnownCause()`.
- Added `Handle()`.
- Fixed detection of existing error-formatting directive failing in special cases.
- Got rid of code generation.
- Added tests, eliminating some issues.

# v0.1.0 (2020-01-28)
- Initial release.
