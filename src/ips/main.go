package main

import (
	"fmt"
	"net"
)

func main() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Print("Error 1: ", err, "\n")
		return
	}
	// handle err
	for _, i := range ifaces {
		fmt.Println("Name:", i.Name)
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print("Error 2: ", err, "\n")
		}
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				fmt.Print("IPNet: ", ip, "\n")
			case *net.IPAddr:
				ip := v.IP
				fmt.Print("IPAddr: ", ip, "\n")
			}
		}
	}
}
