// microservices_bus.go

package main

import (
	"encoding/json"
	"fmt"

	driver "microservices_bus/driver"
	events "microservices_bus/event"
	log "microservices_bus/log"
	"net/http"
	"sync"
)

type Driver interface {
	RegisterService(serviceName, serviceURL string)
	HandleRequest(serviceName, requestData string) (string, error)
	EventListener()
	NotifyEvent(event events.Event)
}

type MicroservicesBus struct {
	driver      Driver
	messageLog  *log.MessageLog
	messageLock sync.Mutex
}

func NewMicroservicesBus(driver Driver, messageLog *log.MessageLog) *MicroservicesBus {

	return &MicroservicesBus{
		driver:     driver,
		messageLog: messageLog,
	}
}

func (mb *MicroservicesBus) RegisterService(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	serviceName, ok := data["name"]
	if !ok {
		http.Error(w, "Service name is required", http.StatusBadRequest)
		return
	}

	serviceURL, ok := data["url"]
	if !ok {
		http.Error(w, "Service URL is required", http.StatusBadRequest)
		return
	}

	mb.driver.RegisterService(serviceName, serviceURL)
	mb.logMessage(serviceName, fmt.Sprintf("Service %s registered successfully", serviceName))

	fmt.Fprintf(w, "Service %s registered successfully", serviceName)
}

func (mb *MicroservicesBus) HandleRequest(w http.ResponseWriter, r *http.Request) {
	serviceName := r.Header.Get("X-Service-Name")
	if serviceName == "" {
		http.Error(w, "X-Service-Name header is required", http.StatusBadRequest)
		return
	}

	// Simulate getting request data
	requestData := "Sample Request Data"

	responseData, err := mb.driver.HandleRequest(serviceName, requestData)
	if err != nil {
		http.Error(w, "Failed to handle request", http.StatusInternalServerError)
		return
	}

	mb.logMessage(serviceName, fmt.Sprintf("Request handled successfully: %s", requestData))

	// Simulate sending response data back to the client
	fmt.Fprintf(w, "Response from %s: %s", serviceName, responseData)
}

func (mb *MicroservicesBus) logMessage(serviceName, data string) {
	mb.messageLock.Lock()
	defer mb.messageLock.Unlock()
	mb.messageLog.LogMessage(serviceName, data)
}

func (mb *MicroservicesBus) GetLogs(w http.ResponseWriter, r *http.Request) {
	mb.messageLock.Lock()
	defer mb.messageLock.Unlock()

	logs := mb.messageLog.GetLogs()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (mb *MicroservicesBus) Start() {
	http.HandleFunc("/register", mb.RegisterService)
	http.HandleFunc("/", mb.HandleRequest)
	http.HandleFunc("/logs", mb.GetLogs)
	http.ListenAndServe(":8080", nil)
}

func main() {
	kafkaDriver := driver.NewKafkaDriver()
	messageLog := log.NewMessageLog()
	microservicesBus := NewMicroservicesBus(kafkaDriver, messageLog)

	go kafkaDriver.EventListener()

	microservicesBus.Start()
}
