package tools

import (
	"math/rand"
	"time"
)

func randint(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(max-min) + min
	return r
}

func Pick(arr []string) string {
	pickItem := arr[randint(0, len(arr))]
	return pickItem
}
