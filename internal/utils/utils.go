package utils

import (
	"math/rand"
)

func RandomizeSlice(slice []string) []string {

	// Fisher-Yates shuffle
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}
