package etapas

import (
	"math/rand"
)

func RandomNumberJugadorE1() int {
	min := 0
	max := 10
	return (rand.Intn(max-min) + min)
}

func RandomNumberLiderE2() int {
	min := 6
	max := 10
	return (rand.Intn(max-min) + min)
}
