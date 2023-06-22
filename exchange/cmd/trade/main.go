package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/alvesrafa/homebroker/exchange/internal/infrastructure/kafka"
	"github.com/alvesrafa/homebroker/exchange/internal/market/dto"
	"github.com/alvesrafa/homebroker/exchange/internal/market/entity"
	"github.com/alvesrafa/homebroker/exchange/internal/market/transformer"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)

	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094", // unix -> /etc/hosts has to have the line host.docker.internal // win -> win/system32/drivers/etc/hosts
		"group.id":          "market",
		"auto.offset.reset": "latest", // earliest start the read from the first message
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"orders"})

	go kafka.Consume(kafkaMsgChan) // specific thread to this function started

	book := entity.NewBook(ordersIn, ordersOut, wg)

	go book.Trade() // specific thread to this function started

	go func() {
		for msg := range kafkaMsgChan {

			wg.Add(1) // start a new task to wait group
			fmt.Println(string(msg.Value))

			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)
			ordersIn <- order

		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(outputJson, []byte("trades"), "output")
	}

}
