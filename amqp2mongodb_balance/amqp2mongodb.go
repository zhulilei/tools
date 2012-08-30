package main

import (
	"flag"
)

var (
	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	mongouri     = flag.String("mongouri", "mongodb://myuser:mypass@localhost:27017/mydatabase", "MONGODB RUI")
	user         = flag.String("user", "admin", "mongodb user")
	password     = flag.String("passwd", "admin", "mongodb password")
	dbname       = flag.String("db", "mydatabase", "mongodb database")
	collection   = flag.String("collection", "metrics", "mongodb collection")
)

const nWorker = 10

type Message struct {
	done    chan int
	content string
}

func task(message_chan chan *Message) {
	consumer := NewConsumer(*uri, *exchange, *exchangeType, *queue, *bindingKey, *consumerTag)
	producer := NewProducer(*mongouri, *dbname, *collection, *user, *password)
	go consumer.read_record(message_chan)
	go producer.insert_record(message_chan)
}
func main() {
	flag.Parse()
	message_chan := make(chan *Message)
	for i := 0; i < nWorker; i++ {
		go task(message_chan)
	}
	select {}
}
