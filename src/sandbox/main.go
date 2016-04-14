package main

import (
	/*
		"fmt"
		"math/rand"
		"time"
	*/

	"fmt"
	"strings"
)

type Object struct {
	Id     int
	Fields map[string]string
}

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

	//match
	obj1 := new(Object)
	obj1.Id = 1
	fields1 := map[string]string{"foo1": "bar1", "foo2": "bar1"}
	obj1.Fields = fields1

	obj2 := new(Object)
	obj2.Id = 2
	fields2 := map[string]string{"foo1": "bar1", "foo2": "bar5"}
	obj2.Fields = fields2

	obj3 := new(Object)
	obj3.Id = 3
	fields3 := map[string]string{"foo1": "bar5", "foo4": "bar1"}
	obj3.Fields = fields3

	obj4 := new(Object)
	obj4.Id = 4
	fields4 := map[string]string{"foo3": "bar5", "foo4": "bar1"}
	obj4.Fields = fields4

	obj5 := new(Object)
	obj5.Id = 5
	fields5 := map[string]string{"foo1": "bar5"}
	obj5.Fields = fields5

	objs := []*Object{obj1, obj2, obj3, obj4, obj5}

	//filtering is done by query parameters on the URI
	filters := map[string][]string{"foo1": []string{"bar1", "bar2", "bar3"},
		"foo2": []string{"bar1", "bar4"}}

	include := true
	for _, obj := range objs {
		for key, values := range filters {
			fmt.Print("Filter Key: ", key, "\n")
			if len(obj.Fields[key]) == 0 {
				fmt.Print("Key ", key, " not found\n")
				include = false
				break
			}
			if !include {
				fmt.Print("Exiting early with no key found\n")
				break
			}

			found := false
			for _, value := range values {
				fmt.Print("Filter Val: ", value, "\n")
				//omit adding to the slice if the key and value doesnt exist
				if strings.Compare(value, obj.Fields[key]) == 0 {
					fmt.Print(value, " = ", obj.Fields[key], "\n")
					found = true //key exists and value exists in the map
					break
				}
			}
			if !found {
				fmt.Print("Exiting early with no value found\n")
				include = false
				break
			}

			fmt.Print("Found: ", found, "\n")
			include = include && found

			if !include {
				fmt.Print("Exiting early with no key found\n")
				break
			}
		}

		if !include {
			fmt.Print("\n\n\n")
			continue
		}

		fmt.Print("ADDING: ", obj.Id, "\n\n\n")
	}

}
