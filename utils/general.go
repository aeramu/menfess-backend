package utils

import "math/rand"

func RandomChance(a, b int) bool {
	i := rand.Int() % b
	if i < a {
		return true
	}
	return false
}
