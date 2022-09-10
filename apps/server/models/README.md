# Server Models

These are broken into two categories, `requests` and `responses`. Requests being the ones coming into pippin and responses being the ones going to clients.

These are overly verbose and some of them are essentially identical, this could probably be re-factored to be better.

The rationale of having a separate unit for each `action` even when parameters are the same is so that they can individually be extended with additional features.

The tests for serializing and deserializing may seem silly, but I've had mistakes before where I mistyped a `json` or `mapstructure` annotation, the tests help prevent that.
