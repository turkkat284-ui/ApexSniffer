package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Logger handles asynchronous logging of AnalyzedPacket data to a file.
type Logger struct {
	file       *os.File
	encoder    *json.Encoder
	logChan    chan AnalyzedPacket
	wg         sync.WaitGroup
	closeChan  chan struct{}
}

// NewLogger creates a new Logger that writes to the specified filePath.
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	l := &Logger{
		file:      file,
		encoder:   json.NewEncoder(file),
		logChan:   make(chan AnalyzedPacket, 10000), // Large buffer to prevent blocking
		closeChan: make(chan struct{}),
	}

	l.wg.Add(1)
	go l.processLogs()

	return l, nil
}

// Log queues a packet for asynchronous logging.
func (l *Logger) Log(packet AnalyzedPacket) {
	select {
	case l.logChan <- packet:
	default:
		// In a real-world scenario, we might want to increment a "dropped logs" counter.
		// For now, we avoid blocking the analyzer if the buffer is full.
	}
}

// processLogs runs in a background goroutine and writes packets to the file.
func (l *Logger) processLogs() {
	defer l.wg.Done()

	for {
		select {
		case packet, ok := <-l.logChan:
			if !ok {
				return
			}
			if err := l.encoder.Encode(packet); err != nil {
				fmt.Fprintf(os.Stderr, "logger error: failed to encode packet: %v\n", err)
			}
		case <-l.closeChan:
			// Drain remaining logs before exiting
			for packet := range l.logChan {
				if err := l.encoder.Encode(packet); err != nil {
					fmt.Fprintf(os.Stderr, "logger error during shutdown: failed to encode packet: %v\n", err)
				}
			}
			return
		}
	}
}

// Close gracefully shuts down the logger, flushing all remaining logs to disk.
func (l *Logger) Close() error {
	close(l.closeChan)
	close(l.logChan)
	l.wg.Wait()
	return l.file.Close()
}
