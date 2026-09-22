# DCK version

This directory contains the construction-kit version of go-fr010. The original Go sources are preserved at their original paths (revision `abdaaecc5e9bc32c4d5a2e4ffd4515ee7abc13a4`), with small asset accessors so both versions use the same embedded resources.

Run the original with `go run ./cmd/fr010` and this version with `go run ./dck/cmd/fr010` from the repository root.

The choreography and assets remain in this repository. Reusable rendering and
effects come from the published `github.com/olivierh59500/democonstructionkit`
module pinned in `go.mod`. Go downloads the dependencies automatically, including
`github.com/olivierh59500/ym-player v1.0.0` for YM playback. Second Reality retains its original ST3 music synchronization.
