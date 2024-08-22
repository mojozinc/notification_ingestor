package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	kafkaBroker   = os.Getenv("KAFKA_BROKER")
	kafkaPassword = os.Getenv("KAFKA_CLIENT_PASSWORDS")
	kafkaTopic    = "alerts"
	producer      *kafka.Producer
)

// In a production setting we can read the kafka config from etcd or consul
func readKafkaConfig() (kafka.ConfigMap, error) {
	// Read the Kafka configuration from a JSON config file
	configFile, err := os.Open("kafka_config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open Kafka config file: %v", err)
	}
	defer configFile.Close()

	var kafkaConfig map[string]string
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&kafkaConfig); err != nil {
		return nil, fmt.Errorf("failed to parse Kafka config file: %v", err)
	}

	// Create a new Kafka producer with the configuration
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"sasl.password":     kafkaPassword,
	}
	for key, value := range kafkaConfig {
		(*configMap)[key] = value
	}
	return *configMap, nil
}

func init() {
	// Create a new Kafka producer
	var err error

	// Read the Kafka configuration
	kafkaConfigMap, err := readKafkaConfig()
	if err != nil {
		log.Fatalf("Failed to read Kafka config: %v", err)
	}

	producer, err = kafka.NewProducer(&kafkaConfigMap)
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
