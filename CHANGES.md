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

# v0.4.0 (2020-02-06)
- Add more info to README.md.
- Renamed `Handle()` to `Recover()`.
- `panik.value` is now `panik.Value`, allowing users to inspect the value.
