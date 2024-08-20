package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	kafkaBroker     = "kafka:9092"
	kafkaGroupID    = "alert-consumer-group"
	kafkaTopic      = "alerts"
	mongoURI        = "mongodb://mongo:27017"
	mongoDBName     = "alert_db"
	mongoCollection = "alerts"
)

func main() {
	// Kafka consumer configuration
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          kafkaGroupID,
		"auto.offset.reset": "earliest",
	}

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Subscribe to Kafka topic
	err = consumer.Subscribe(kafkaTopic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	// MongoDB client configuration
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(mongoDBName).Collection(mongoCollection)

	fmt.Println("Consumer started...")

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

		fmt.Printf("Received message: %+v\n", alert)

		// Insert into MongoDB
		_, err = collection.InsertOne(context.TODO(), alert)
		if err != nil {
			log.Printf("Error inserting into MongoDB: %v", err)
			continue
		}

		fmt.Println("Inserted alert into MongoDB")
	}
}
