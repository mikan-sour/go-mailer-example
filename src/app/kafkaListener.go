package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jedzeins/go-mailer/src/config"

	"github.com/Shopify/sarama"
)

func (app *AppImpl) RunKafkaListener(config *config.Config) {
	worker, err := app.KafkaListenerService.SetupConsumer(app.KafkaListenerService.Ports)
	if err != nil {
		app.ErrorLog.Fatalf("Error RunKafkaListener: %s", err.Error())
		panic(err)
	}

	consumer, err := worker.ConsumePartition(config.KAFKA_LISTEN_TOPIC, 0, sarama.OffsetNewest)
	if err != nil {
		app.ErrorLog.Fatalf("Error ConsumePartition: %s", err.Error())
		panic(err)
	}

	app.InfoLog.Println("Consumer Partition is being listened on")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	// Count how many message processed

	// Get signal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				app.ErrorLog.Fatalf("Error in consumer select case: %s", err.Error())
			case msg := <-consumer.Messages():
				app.InfoLog.Printf("Received new message at %s", msg.Timestamp)
				parsedMessage, err := app.KafkaListenerService.ParseMessage(msg.Value)
				if err != nil {
					app.ErrorLog.Fatalf("Error in ParseMessage select case: %s", err.Error())
				}
				app.SendEmail(*parsedMessage)
			case <-sigchan:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	app.InfoLog.Println("Listener is done")

	if err := worker.Close(); err != nil {
		app.ErrorLog.Fatalf("Error closing the worker: %s", err.Error())
		panic(err)
	}

}
