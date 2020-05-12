# v0.1.0 (2020-01-28)
- Initial release.

# v0.2.0 (2020-01-29)
- Added `OnError()`.
- Added `Panic()`.
- Added `IsKnownCause()`.
- Added `Handle()`.
- Fixed detection of existing error-formatting directive failing in special cases.
- Got rid of code generation.
- Added tests, eliminating some issues.

# v0.3.0 (2020-01-31)
- Major rewrite.

# v0.4.0 (2020-02-07)
- Added more info to README.md.
- Removed `Handle()`.
- `panik.value` is now `panik.Value`, allowing users to inspect the value.
- Changed signature of `Panic()` to be consistent with `panic()`.
- Simplified API.

# v0.4.1 (2020-04-21)
- Added `RecoverTraceFunc()` and `ExitTraceFunc()`.

# v0.5.0 (2020-05-12)
- Removed `RecoverTrace()` and `ExitTrace()`.
- Added `ToErrorWithTrace()`.
