package main

import (
	"fmt"
	"os"
)

/*
	"fmt"
	"math/rand"
	"time"
*/

func main() {
	/*
		test := []string{"string1", "string2", "string3"}

		rand.Seed(time.Now().UnixNano())

		for i := 0; i < 10; i++ {
			random := rand.Intn(len(test))
			fmt.Println(test[random])
		}
	*/

	/*
		test := "ping 1"
		iindex := strings.Index(test, " ")
		value := test[iindex+1:]
		fmt.Println("value: ", value)
	*/

	for i := 0; i < len(os.Args); i++ {
		fmt.Println("Arg", (i + 1), ":", os.Args[i])
	}

	fmt.Println("hello")
	fmt.Println("how")
	fmt.Println("are")
	fmt.Println("you")
}
