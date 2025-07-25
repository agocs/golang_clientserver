package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/agocs/golang_clientserver/payload"
)

// ThrottledReader wraps an io.Reader and throttles reading speed
type ThrottledReader struct {
	reader       io.Reader
	bytesPerRead int
	delay        time.Duration
}

// NewThrottledReader creates a new ThrottledReader
// totalBytes: total bytes to read
// duration: desired total read duration
func NewThrottledReader(reader io.Reader, totalBytes int, duration time.Duration) *ThrottledReader {
	// Calculate how many reads we want to do in the given duration
	const readsCount = 100 // divide reading into 100 chunks
	bytesPerRead := totalBytes / readsCount
	if bytesPerRead < 1 {
		bytesPerRead = 1
	}

	// Calculate delay between reads
	delay := duration / time.Duration(readsCount)

	return &ThrottledReader{
		reader:       reader,
		bytesPerRead: bytesPerRead,
		delay:        delay,
	}
}

// Read implements io.Reader interface with throttling
func (t *ThrottledReader) Read(p []byte) (int, error) {
	// Limit how much we read in one call
	toRead := len(p)
	if toRead > t.bytesPerRead {
		toRead = t.bytesPerRead
	}

	// Read from the underlying reader
	n, err := t.reader.Read(p[:toRead])

	// Delay after reading
	time.Sleep(t.delay)

	return n, err
}

// generateLargeRandomString creates a random string of the specified size in MB
func generateLargeRandomString(sizeMB int) string {
	// Calculate bytes for the given MB
	sizeBytes := sizeMB * 1024 * 1024

	// Character set to use for the random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice of the required size
	b := make([]byte, sizeBytes)

	// Fill it with random characters
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func main() {
	// Parse command line arguments
	throttleRequest := false
	for _, arg := range os.Args[1:] {
		if arg == "--throttled" {
			throttleRequest = true
		}
	}

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go doRequests(wg, throttleRequest)
	}
	wg.Wait()
}

func doRequests(wg *sync.WaitGroup, throttleRequest bool) {
	defer wg.Done()
	for {
		err := makeRequest(throttleRequest)
		if err != nil {
			return
		}
	}
}

func makeRequest(throttleRequest bool) error {
	largeContents := generateLargeRandomString(10)

	payload := payload.Payload{
		SentTime: time.Now(),
		Contents: largeContents,
	}

	payloadJSON, _ := json.Marshal(payload)

	// Log start time
	startTime := time.Now()
	requestTypeString := "normal"
	if throttleRequest{
		requestTypeString = "throttled"
	}
	log.Printf("Request started at: %s using %s mode", startTime.Format("15:04:05.000000000"), requestTypeString)

	// Create a throttled reader for the payload data
	// This will make the data transfer take approximately 1 second
	normalReader := bytes.NewReader(payloadJSON)
	throttledReader := NewThrottledReader(normalReader, len(payloadJSON), 1*time.Second)

	var readerToUse io.Reader
	if throttleRequest {
		readerToUse = throttledReader
	} else {
		readerToUse = normalReader
	}

	// Send the request with the selected reader
	resp, err := http.Post("http://localhost:8080", "application/json", readerToUse)
	if err != nil {
		return err
	}

	// Log elapsed time
	elapsed := time.Since(startTime)

	log.Printf("Response status: %s", resp.Status)
	log.Printf("Sent payload size: %.2f MB", float64(len(payloadJSON))/(1024*1024))
	log.Printf("Request took: %v", elapsed)
	toSleep := rand.Intn(10)
	log.Printf("Sleeping for %d", toSleep)
	time.Sleep(time.Second * time.Duration(toSleep))
	return nil
}
