package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

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

	filename := path + "/mydata.txt"

	go func() {
		count := 0

		inFile, err := os.Open(filename)
		if err == nil {
			defer inFile.Close()

			scanner := bufio.NewScanner(inFile)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				count++
			}
		}

		count++ //so we can start with 1 and continue where we left off

		for {
			outFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			failOnError(err, "Failed to open file for append")
			defer outFile.Close()

			_, err = outFile.WriteString("Creating my data #" + strconv.Itoa(count) + "\n")
			failOnError(err, "Failed to append file")

			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(7) + 1
			time.Sleep(time.Duration(random) * time.Second)

			count++
		}

	}()

	panic(http.ListenAndServe(fmt.Sprintf(":%v", port), http.FileServer(http.Dir(path))))
}
