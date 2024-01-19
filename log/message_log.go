// log/message_log.go

package log

import (
	"fmt"
	"log"
	"sync"
)

// LogEntry represents a log entry.
type LogEntry struct {
	ServiceName string `json:"serviceName"`
	Data        string `json:"data"`
}

// MessageLog represents a message log.
type MessageLog struct {
	entries []LogEntry
	mutex   sync.Mutex
}

// NewMessageLog creates a new MessageLog instance.
func NewMessageLog() *MessageLog {
	log.Default().Println("Initialising Message Logger")
	return &MessageLog{
		entries: make([]LogEntry, 0),
	}
}

// LogMessage logs a message.
func (ml *MessageLog) LogMessage(serviceName, data string) {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()
	entry := LogEntry{
		ServiceName: serviceName,
		Data:        data,
	}
	fmt.Printf("ServiceName ")
	ml.entries = append(ml.entries, entry)
}

// GetLogs returns the log entries.
func (ml *MessageLog) GetLogs() []LogEntry {
	ml.mutex.Lock()
	defer ml.mutex.Unlock()
	return ml.entries
}
