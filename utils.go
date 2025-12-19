package main

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func absMax(a, b int) int {
	aa := absInt(a)
	bb := absInt(b)
	if aa > bb {
		return aa
	}
	return bb
}
