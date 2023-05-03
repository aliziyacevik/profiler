package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// userAppHandler simulates the user's app running on localhost:50133.
func userAppHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from user's app!"))
}

func TestProfilerUnbufferedWrite(t *testing.T) {
	// Set up the user's app server.
	userAppServer := httptest.NewServer(http.HandlerFunc(userAppHandler))
	defer userAppServer.Close()

	// Set up the Profiler configuration.
	config := Config{
		Port:         50132,
		OutFile:      "test_profile.json",
		OutFormat:    "JSON",
		TargetServer: userAppServer.URL,
		BufferedWrite: false,
	}

	// Remove the test output file after the test is done.
	//	defer os.Remove(config.OutFile)

	// Start the Profiler server.
	go func() {
		startProfilerServer(config)
	}()

	// Wait for the Profiler server to start.
	time.Sleep(1 * time.Second)

	// Send a request to the Profiler server.
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/test", config.Port))
	if err != nil {
		t.Fatalf("Failed to send request to Profiler server: %v", err)
	}
	defer resp.Body.Close()

	// Check if the Profiler server successfully proxies the request to the user's app server.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := "Hello from user's app!"
	if string(body) != expectedBody {
		t.Errorf("Expected response body: %s, got: %s", expectedBody, string(body))
	}

	// Check if the Profiler server saves the request information to the output file.
	_, err = os.Stat(config.OutFile)
	if os.IsNotExist(err) {
		t.Errorf("Output file not created: %s", config.OutFile)
	}
}
/*
func TestProfilerBufferedWrites(t *testing.T) {
	// Set up the user's app server.
	userAppServer := httptest.NewServer(http.HandlerFunc(userAppHandler))
	defer userAppServer.Close()

	// Set up the Profiler configuration.
	config := Config{
		Port:         50132,
		OutFile:      "test_profile_buffered.json",
		OutFormat:    "JSON",
		TargetServer: userAppServer.URL,
		BufferedWrite: true,
	}

	// Remove the test output file after the test is done.
	defer os.Remove(config.OutFile)

	// Start the Profiler server.
	go func() {
		startProfilerServer(config)
	}()

	// Wait for the Profiler server to start.
	time.Sleep(1 * time.Second)
	bufferSize := 100 
	// Send multiple requests to the Profiler server.
	for i := 0; i < bufferSize+1; i++ {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/test", config.Port))
		if err != nil {
			t.Fatalf("Failed to send request to Profiler server: %v", err)
		}
		defer resp.Body.Close()

		// Check if the Profiler server successfully proxies the request to the user's app server.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedBody := "Hello from user's app!"
		if string(body) != expectedBody {
			t.Errorf("Expected response body: %s, got: %s", expectedBody, string(body))
		}
	}

	// Wait for the buffer to be processed.
	time.Sleep(1 * time.Second)

	// Check if the Profiler server saves the request information to the output file.
	_, err := os.Stat(config.OutFile)
	if os.IsNotExist(err) {
		t.Errorf("Output file not created: %s", config.OutFile)
	}
}
*/
