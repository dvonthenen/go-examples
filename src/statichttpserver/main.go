package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	//define flags
	var path string
	var port int
	flag.StringVar(&path, "path", ".", "directory from which to serve files from")
	flag.IntVar(&port, "port", 8080, "the port in which to bind to")
	//parse
	flag.Parse()

	fmt.Println("path = ", path)
	fmt.Println("port = ", port)

	panic(http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir(path))))
}
