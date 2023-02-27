package common

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkNaive(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	// generate rand lengths

	N := 10000
	nums := make([]int, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Intn(N) + 1
	}
	b.Run(
		fmt.Sprintf("input_size_%d", N),
		func(b *testing.B) {
			j := 0
			for i := 0; i < b.N; i++ {
				if j > N-1 {
					j = 0
				}
				rand.Intn(nums[j])
				j++
			}
		})

}

func BenchmarkFloat(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	N := 10000
	nums := make([]int, N)
	for i := 0; i < N; i++ {
		nums[i] = rand.Intn(N) + 1
	}

	b.Run(
		fmt.Sprintf("input_size_%d", N),
		func(b *testing.B) {
			j := 0
			for i := 0; i < b.N; i++ {
				if j > N-1 {
					j = 0
				}
				math.Round(GetRandFloat() * float64(nums[j]))
				j++
			}
		})
}
