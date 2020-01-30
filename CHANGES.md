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

# v0.3.0 (TBD)
- Renamed `WriteTrace()` to `RecoverTraceTo()`.
- Added `RecoverTrace()`.
- `ToError()` will no longer contain a `*panik.knownCause` in the error-chain.
- Removed `ToCustomError()`.
- Renamed `IsKnownCause()` to `HasKnownCause()`.
