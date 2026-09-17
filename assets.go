package fr010

import _ "embed"

//go:embed data/kasparov.bin
var kasparovData []byte

//go:embed data/bug.bin
var bugData []byte

//go:embed data/stadt.bin
var stadtData []byte

//go:embed "data/Stormlord 1 - title.ym"
var ymData []byte

//go:embed data/secrcode.ttf
var secrcodeData []byte
