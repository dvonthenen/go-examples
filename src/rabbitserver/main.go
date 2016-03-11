package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type session struct {
	x map[string][]string
	sync.Mutex
}

func (s *session) add(key string, str string) {
	s.Lock()
	s.x[key] = append(s.x[key], str)
	s.Unlock()
}

func (s *session) print() string {
	s.Lock()
	str := "<html><head><title>Output</title></head><body>"

	for k, v := range s.x {
		str += "Client "
		str += k
		str += "<br />"
		for _, m := range v {
			str += "&nbsp;&nbsp;&nbsp;"
			str += m
			str += "<br />"
		}
	}

	str += "</body></html>"
	s.Unlock()

	return str
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello.rabbitserver.queue", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	//keep history
	s := new(session)
	s.x = make(map[string][]string)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("<- Received %s from %s", d.Body, d.CorrelationId)

			s.add(d.CorrelationId, string(d.Body))

			str := string(d.Body)
			if !strings.HasPrefix(str, "ping ") {
				log.Printf("Invalid string format")
				continue
			}

			newstr := strings.Replace(str, "ping ", "", -1)

			msg := "pong " + newstr
			log.Printf("-> Sending %s to %s", msg, d.CorrelationId)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(msg),
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveRest(w, r, s)
	})
	http.ListenAndServe(":9000", nil)

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func serveRest(w http.ResponseWriter, r *http.Request, s *session) {
	fmt.Fprintf(w, s.print())
}
