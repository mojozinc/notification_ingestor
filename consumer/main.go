package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	kafkaBroker   = os.Getenv("KAFKA_BROKER")
	kafkaPassword = os.Getenv("KAFKA_CLIENT_PASSWORDS")
	kafkaGroupID  = "alert-consumer-group"
	kafkaTopic    = "alerts"
)

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
		"group.id":          kafkaGroupID,
		"auto.offset.reset": "earliest",
	}
	for key, value := range kafkaConfig {
		(*configMap)[key] = value
	}
	return *configMap, nil
}

func main() {
	// Kafka consumer configuration

	kafkaConfig, err := readKafkaConfig()
	if err != nil {
		log.Fatalf("Failed to read Kafka config: %v", err)
	}
	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(&kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to Kafka topic
	err = consumer.Subscribe(kafkaTopic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// The writer is an abstraction for persisting messages
	// We can store message in a NoSQL database like ElasticSearch, time-series database like InfluxDB
	// We can think about further optimization for read patterns like, fetching alerts by column, table, time range etc.

	writer := DiskWriter{}
	log.Println("Consumer started...")

	for {
		// Poll for new messages
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Parse message
		var alert map[string]interface{}
		if err := json.Unmarshal(msg.Value, &alert); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		log.Printf("Received message: %+v\n", alert)
		writer.WriteMessage(alert)
		log.Println("Message persisted in storage")
	}
}
