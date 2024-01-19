// kafka_driver.go

package driver

import (
	"fmt"
	"log"
	event "microservices_bus/event"
	"os"
	"os/signal"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaDriver struct {
	producer     *kafka.Producer
	consumer     *kafka.Consumer
	eventChannel chan event.Event
	services     map[string]string
	mutex        sync.Mutex
}

func NewKafkaDriver() *KafkaDriver {
	log.Default().Println("Initialising Kafka Driver")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9093", "api.version.request": true})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	} else {
		log.Default().Println("Kafka producer Initialised")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   "kafka:9093",
		"group.id":            "my-group",
		"auto.offset.reset":   "earliest",
		"api.version.request": true,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}else {
		log.Default().Println("Kafka consumer Initialised")
	}

	return &KafkaDriver{
		producer:     p,
		consumer:     c,
		eventChannel: make(chan event.Event, 10),
		services:     make(map[string]string),
	}
}

func (d *KafkaDriver) RegisterService(serviceName, serviceURL string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.services[serviceName] = serviceURL
}

func (d *KafkaDriver) HandleRequest(serviceName, requestData string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, ok := d.services[serviceName]
	if !ok {
		return "", fmt.Errorf("Service not found")
	}

	// Simulate sending request to the target service
	responseData := fmt.Sprintf("Response from %s: %s", serviceName, requestData)

	// Notify the event
	d.eventChannel <- event.Event{ServiceName: serviceName, Data: "Request completed successfully"}

	return responseData, nil
}

func (d *KafkaDriver) EventListener() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	defer close(sigchan)

	for {
		select {
		case ev := <-d.eventChannel:
			// Simulate sending event to Kafka topic
			message := &kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &ev.ServiceName, Partition: kafka.PartitionAny},
				Value:          []byte(ev.Data),
			}
			if err := d.producer.Produce(message, nil); err != nil {
				log.Printf("Failed to produce event: %v", err)
			}
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating", sig)
			return
		}
	}
}

func (d *KafkaDriver) NotifyEvent(event event.Event) {
	// Notify the event
	d.eventChannel <- event
}

// Additional method to get service URL based on the Kafka topic
func (d *KafkaDriver) GetServiceURL(serviceName string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	serviceURL, ok := d.services[serviceName]
	if !ok {
		return "", fmt.Errorf("Service not found")
	}
	return serviceURL, nil
}
