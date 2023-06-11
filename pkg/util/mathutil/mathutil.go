package mathutil

import (
	"golang.org/x/exp/constraints"
	"math"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Pad(totalSize, chunkSize int) int {
	return int(math.Ceil(float64(totalSize)/float64(chunkSize)) * float64(chunkSize))
}
