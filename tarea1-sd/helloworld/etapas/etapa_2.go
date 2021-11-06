package etapas

import (
	"math/rand"
)

func RandomNumberStageTwo() int {

	min := 1
	max := 4
	return (rand.Intn(max-min) + min)

}
