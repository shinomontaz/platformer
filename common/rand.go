package common

import (
	"math/rand"
	"time"
)

var (
	randInt   []int
	randFloat []float64
	rintc     int
	rfloatc   int
)

const MAXRANDNUM = 100000

func init() {
	randInt = make([]int, MAXRANDNUM)
	randFloat = make([]float64, MAXRANDNUM)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < MAXRANDNUM; i++ {
		randInt[i] = rand.Intn(10)
		randFloat[i] = rand.Float64()
	}
}

func GetRandInt() int {
	if rintc >= MAXRANDNUM-1 {
		rintc = -1
	}
	rintc++

	return randInt[rintc]
}

func GetRandFloat() float64 {
	if rfloatc >= MAXRANDNUM-1 {
		rfloatc = -1
	}
	rfloatc++

	return randFloat[rfloatc]
}
