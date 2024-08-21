package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	kafkaBroker = os.Getenv("KAFKA_BROKER")
	kafkaTopic  = "alerts"
	producer    *kafka.Producer
)

func init() {
	// Create a new Kafka producer
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"security.protocol": "SASL_PLAINTEXT",
		"sasl.mechanism":    "PLAIN",
		"sasl.username":     "user1",
		"sasl.password":     "U15z0kC9yC",
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
}

func handleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	message, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to serialize payload", http.StatusInternalServerError)
		return
	}

	log.Printf("Received notification: %s", message)

	// Produce the message to Kafka
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)

	if err != nil {
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		http.Error(w, "Failed to deliver message", http.StatusInternalServerError)
		return
	}

	log.Printf("Message delivered to %v", m.TopicPartition)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification received"))
}

func main() {
	http.HandleFunc("/notifications", handleNotification)
	log.Fatal(http.ListenAndServe(":80", nil))
}
