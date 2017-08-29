package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	dataPath   string
	listenPort string
)

const (
	defaultPath = "."
	defaultPort = "8080"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func env(key, defaultValue string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultValue
	}
	return
}

func main() {
	dataPath = env("DATA_PATH", defaultPath)
	listenPort = env("LISTEN_PORT", defaultPort)

	fmt.Println("path = ", dataPath)
	fmt.Println("port = ", listenPort)

	filename := dataPath + "/mydata.txt"

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

	panic(http.ListenAndServe(fmt.Sprintf(":%s", listenPort), http.FileServer(http.Dir(dataPath))))
}
