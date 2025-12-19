package main

var randSeed = 31337303

func randInt() int {
	randSeed = (randSeed*15625 + 1) & 32767
	return randSeed
}

func randFloat() float32 {
	return float32(randInt()) / 32767.0
}

func randByte() byte {
	return byte(randInt())
}
