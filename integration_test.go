// integration_test.go

package main

import (
	"bytes"
	"encoding/json"
	driver "microservices_bus/driver"
	log "microservices_bus/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMicroservicesBusIntegration(t *testing.T) {
	kafkaDriver := driver.NewKafkaDriver()
	messageLog := log.NewMessageLog()
	microservicesBus := NewMicroservicesBus(kafkaDriver, messageLog)

	// Start the microservices bus in a goroutine
	go microservicesBus.Start()

	// Simulate registering a service
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(`{"name":"service1","url":"http://localhost:8081"}`))
	microservicesBus.RegisterService(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200, got %d", resp.Code)
	}

	// Simulate handling a request
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Service-Name", "service1")
	microservicesBus.HandleRequest(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200, got %d", resp.Code)
	}

	// Simulate querying logs
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/logs", nil)
	microservicesBus.GetLogs(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected HTTP status 200, got %d", resp.Code)
	}

	// Unmarshal response body and check the logs
	var logs []log.LogEntry
	if err := json.Unmarshal(resp.Body.Bytes(), &logs); err != nil {
		t.Errorf("Error decoding JSON response: %v", err)
	}

	if len(logs) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(logs))
	}
}
