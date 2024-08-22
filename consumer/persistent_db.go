package main

import (
	"encoding/json"
	"os"
)

// MessageWriter defines the interface for writing messages into a database.
type MessageWriter interface {
	WriteMessage(map[string]interface{}) error
}

// MongoDBWriter is an implementation of MessageWriter that connects to MongoDB.
type MongoDBWriter struct {
	// Add MongoDB connection details here
}

// WriteMessage writes the given message to MongoDB.
func (w *MongoDBWriter) WriteMessage(message map[string]interface{}) error {
	// Implement the logic to write the message to MongoDB
	return nil
}

// DiskWriter is an implementation of MessageWriter that writes the data to disk.
type DiskWriter struct {
	// Add disk storage details here
}

// WriteMessage writes the given message to disk.
func (w *DiskWriter) WriteMessage(message map[string]interface{}) error {
	// Implement the logic to write the message to disk
	// For example, you can write the message to a file
	file, err := os.OpenFile("/alerts_data.ndjson", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert the message to a string representation
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	messageStr := string(messageBytes)

	// Write the message to the file
	_, err = file.WriteString(messageStr)
	if err != nil {
		return err
	}

	return nil
}
