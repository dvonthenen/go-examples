package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	test := []string{"string1", "string2", "string3"}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		random := rand.Intn(len(test))
		fmt.Println(test[random])
	}
}
