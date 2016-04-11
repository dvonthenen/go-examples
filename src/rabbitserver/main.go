package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type session struct {
	address string
	x       map[string][]string
	sync.Mutex
}

func (s *session) add(key string, str string) {
	s.Lock()
	s.x[key] = append(s.x[key], str)
	s.Unlock()
}

func (s *session) printMessages() string {
	s.Lock()
	str := "<html><head><title>Output</title><meta http-equiv=\"refresh\" content=\"2\" /></head><body>"

	var keys []string
	for tk := range s.x {
		keys = append(keys, tk)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str += "Client "
		str += k
		str += "<br />"
		for _, m := range s.x[k] {
			str += "&nbsp;&nbsp;&nbsp;"
			str += m
			str += "<br />"
		}
	}

	/*
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
	*/

	str += "</body></html>"
	s.Unlock()

	return str
}

func (s *session) printDatabase() string {
	s.Lock()
	str := "<html><head><title>Output</title><meta http-equiv=\"refresh\" content=\"2\" /></head><body>"

	db, err := sql.Open("postgres", "host="+s.address+" user=dev password=vmware dbname=demo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//SELECT
	rows, err := db.Query("SELECT id, client,	received_message,	sent_message FROM message ORDER BY client DESC, id ASC")
	if err != nil {
		panic(err)
	}

	lastclient := ""

	for rows.Next() {
		var id int
		var client string
		var receivedMessage string
		var sentMessage string
		err = rows.Scan(&id, &client, &receivedMessage, &sentMessage)
		if err != nil {
			continue
		}

		fmt.Println("id:", id, " client:", client, " receivedMessage:", receivedMessage,
			" sentMessage:", sentMessage)

		if !strings.EqualFold(client, lastclient) {
			str += "Client "
			str += client
			str += "<br />"
		}
		lastclient = client

		str += "&nbsp;&nbsp;&nbsp;"
		str += receivedMessage
		str += " <---> "
		str += sentMessage
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

func addTransactionToDb(address string, client string, received string, response string) {
	db, err := sql.Open("postgres", "host="+address+" user=dev password=vmware dbname=demo sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("Received: %s, Response: %s\n", received, response)

	//INSERT
	var messageid int
	err = db.QueryRow("INSERT INTO message (client, received_message, sent_message) VALUES ($1, $2, $3) RETURNING id",
		client, received, response).Scan(&messageid)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID: ", messageid)
}

func main() {
	//define flags
	var port int
	flag.IntVar(&port, "port", 9000, "the port in which to bind the HTTP server to")
	var address string
	flag.StringVar(&address, "address", "127.0.0.1", "the rabbit server in which to bind to")
	var rabbitservice string
	flag.StringVar(&rabbitservice, "rabbitservice", "", "the rabbit service to autodiscover")
	var postgresaddress string
	flag.StringVar(&postgresaddress, "postgresaddress", "", "the postgres server to connect to")
	//parse
	flag.Parse()

	connstr := "amqp://guest:guest@" + address + ":5672/"
	if len(rabbitservice) > 0 {
		connstr = generateRabbitURI(rabbitservice)
	}

	conn, err := amqp.Dial(connstr)
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
	s.address = postgresaddress

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
			istr, _ := strconv.Atoi(newstr)

			msg := "pong " + newstr
			switch istr % 6 {
			case 0:
				msg = "pong " + newstr
			case 1:
				msg = "goodbye " + newstr
			case 2:
				msg = "adios " + newstr
			case 3:
				msg = "zaijian " + newstr
			case 4:
				msg = "arrivederci " + newstr
			case 5:
				msg = "sayonara " + newstr
			}

			log.Printf("-> Sending %s to %s", msg, d.CorrelationId)

			addTransactionToDb(postgresaddress, d.CorrelationId, str, msg)

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

	//http server to display history
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		printMessages(w, r, s)
	})
	http.HandleFunc("/database", func(w http.ResponseWriter, r *http.Request) {
		printDatabase(w, r, s)
	})
	http.ListenAndServe(":"+strconv.Itoa(port), nil)

	<-forever
}

func printMessages(w http.ResponseWriter, r *http.Request, s *session) {
	fmt.Fprintf(w, s.printMessages())
}

func printDatabase(w http.ResponseWriter, r *http.Request, s *session) {
	fmt.Fprintf(w, s.printDatabase())
}
