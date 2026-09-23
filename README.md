
## Optional DCK version

The original implementation remains at its original paths. Run it with `go run ./cmd/fr010`.

The construction-kit version is in [dck/](dck/README.md). Run `go run ./dck/cmd/fr010` from this directory. Both versions share the original assets.

The DCK version opens its soundtrack with `sound.Open(filename, data, options)`.
DCK selects the decoder internally and preserves the original half-volume,
filtered stereo PCM. No YM decoder or PCM conversion is duplicated in `dck/`.
