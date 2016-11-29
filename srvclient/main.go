package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	fmt.Println("USING AUTODISCOVERY")
	_, srvs, err := net.LookupSRV("restapi", "tcp", "marathon.mesos")
	if err != nil {
		panic(err)
	}
	if len(srvs) == 0 {
		fmt.Println("got no record")
	}
	for _, srv := range srvs {
		fmt.Println("Discovered service:", srv.Target, "port", srv.Port)
	}

	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(len(srvs))
	url := "http://" + srvs[random].Target + ":" + strconv.Itoa(int(srvs[random].Port)) + "/user"
	fmt.Print(url + "\n")
}
