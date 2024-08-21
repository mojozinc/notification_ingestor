// handle_notification_test.go
package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// // Mock Kafka Producer
// type MockProducer struct {
//     mock.Mock
// }

// func (m *MockProducer) Produce(msg *kafka.Message, deliveryChan chan kafka.Event) error {
//     args := m.Called(msg, deliveryChan)
//     if args.Get(0) != nil {
//         return args.Get(0).(error)
//     }
//     return nil
// }

// Helper function to read test data
func readTestData(filename string) ([]byte, error) {
	cwd, _ := os.Getwd()
	println("cwd = ", cwd)
	path := filepath.Join(cwd, "test_data", filename)
	return ioutil.ReadFile(path)
}

func TestHandleNotification_Success(t *testing.T) {

	// Load valid payload from file
	payloadBytes, err := readTestData("valid_payload.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBuffer(payloadBytes))
	w := httptest.NewRecorder()

	// Call the function
	handleNotification(w, req)

	// Check the status code
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200 but got %d", res.StatusCode)
	}

	// Check response body
	body, _ := ioutil.ReadAll(res.Body)
	expectedBody := "Notification received"
	if string(body) != expectedBody {
		t.Errorf("Expected response body %q but got %q", expectedBody, string(body))
	}
}

// func TestHandleNotification_InvalidMethod(t *testing.T) {
// 	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
// 	w := httptest.NewRecorder()

// 	handleNotification(w, req)

// 	res := w.Result()
// 	if res.StatusCode != http.StatusMethodNotAllowed {
// 		t.Errorf("Expected status code 405 but got %d", res.StatusCode)
// 	}
// }

// func TestHandleNotification_InvalidPayload(t *testing.T) {
//     // Load invalid payload from file
//     payloadBytes, err := readTestData("invalid_payload.json")
//     if err != nil {
//         t.Fatalf("Failed to read test data: %v", err)
//     }

//     req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBuffer(payloadBytes))
//     w := httptest.NewRecorder()

//     handleNotification(w, req)

//     res := w.Result()
//     if res.StatusCode != http.StatusBadRequest {
//         t.Errorf("Expected status code 400 but got %d", res.StatusCode)
//     }
// }
