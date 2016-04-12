package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
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

func generateRabbitURI(service string) string {
	fmt.Println("USING AUTODISCOVERY")
	_, srvs, err := net.LookupSRV(service, "tcp", "marathon.mesos")
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

	return "amqp://guest:guest@" + srvs[random].Target + ":5672/"
}

func main() {
	//define flags
	var port int
	flag.IntVar(&port, "port", 8000, "the port in which to bind the HTTP server to")
	var address string
	flag.StringVar(&address, "address", "127.0.0.1", "the rabbit server in which to bind to")
	var service string
	flag.StringVar(&service, "service", "", "the rabbit service to autodiscover")
	//parse
	flag.Parse()

	connstr := "amqp://guest:guest@" + address + ":5672/"
	if len(service) > 0 {
		connstr = generateRabbitURI(service)
	}

	conn, err := amqp.Dial(connstr)
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
			switch count % 5 {
			case 0:
				msg = "ping " + strconv.Itoa(count)
			case 1:
				msg = "hello " + strconv.Itoa(count)
			case 2:
				msg = "hola " + strconv.Itoa(count)
			case 3:
				msg = "ciao " + strconv.Itoa(count)
			case 4:
				msg = "konnichiwa " + strconv.Itoa(count)
			}
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
			random := rand.Intn(4) + 3
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
