package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type session struct {
	msgs []string
	sync.Mutex
}

func (s *session) add(str string) {
	s.Lock()
	s.msgs = append(s.msgs, str)
	s.Unlock()
}

func (s *session) print() string {
	s.Lock()
	str := "<html><head><title>Output</title><meta http-equiv=\"refresh\" content=\"2\" /></head><body>"

	str += "Messages<br />"
	for _, m := range s.msgs {
		str += "&nbsp;&nbsp;&nbsp;"
		str += m
		str += "<br />"
	}

	str += "</body></html>"
	s.Unlock()

	return str
}

func main() {
	//define flags
	var port int
	flag.IntVar(&port, "port", 9001, "the port in which to bind the HTTP server to")
	var address string
	flag.StringVar(&address, "address", "127.0.0.1", "the rabbit server in which to bind to")
	//parse
	flag.Parse()

	conn, err := amqp.Dial("amqp://guest:guest@" + address + ":5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	u4, err := uuid.NewV4()
	failOnError(err, "Failed to generate a uuid")
	log.Printf("Client %s", u4.String())

	s := new(session)
	s.msgs = make([]string, 0)

	forever := make(chan bool)

	go func() {
		count := 1

		for {
			msg := "ping " + strconv.Itoa(count)
			log.Printf("-> Sending %s", msg)

			err = ch.Publish(
				"", // exchange
				"hello.rabbitserver.queue", // routing key
				false, // mandatory
				false, // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: u4.String(),
					ReplyTo:       q.Name,
					Body:          []byte(msg),
				})
			failOnError(err, "Failed to publish a message")

			for d := range msgs {
				if u4.String() == d.CorrelationId {
					log.Printf("<- Received %s", d.Body)
					s.add(string(d.Body))
					break
				}
			}

			rand.Seed(time.Now().UnixNano())
			random := rand.Intn(4) + 1
			time.Sleep(time.Duration(random) * time.Second)

			count++
		}
	}()

	//http server to display history
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveRest(w, r, s)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)

	<-forever
}

func serveRest(w http.ResponseWriter, r *http.Request, s *session) {
	fmt.Fprintf(w, s.print())
}
