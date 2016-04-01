package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	cname, srvs, err := net.LookupSRV("postgres", "tcp", "marathon.mesos")
	if err != nil {
		panic(err)
	}
	if len(srvs) == 0 {
		fmt.Println("got no record")
	}
	if !strings.HasSuffix(cname, "marathon.mesos") {
		fmt.Println("got", cname, "want marathon.mesos")
	}
	for _, srv := range srvs {
		if !strings.HasSuffix(srv.Target, "marathon.mesos") {
			fmt.Println("got", srv, "want a record containing marathon.mesos")
		}
	}
}
