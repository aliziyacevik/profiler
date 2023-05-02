package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         int
	OutFile      string
	OutFormat    string
	TargetServer string
}

type RequestInfo struct {
	Method   string            `json:"method"`
	URL      string            `json:"url"`
	Browser  string            `json:"browser"`
	Header   map[string]string `json:"header"`
	IP       string            `json:"ip"`
	Received string            `json:"received"`
}

func loadConfig(configFile string) (Config, error) {
	err := godotenv.Load(configFile)
	if err != nil {
		return Config{}, err
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:      port,
		OutFile:   os.Getenv("OUT_FILE"),
		OutFormat: os.Getenv("OUT_FORMAT"),
	}, nil
}

func writeRequestInfoToFile(requestInfo *RequestInfo, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var requestInfos []RequestInfo
	if len(data) > 0 {
		err = json.Unmarshal(data, &requestInfos)
		if err != nil {
			return err
		}
	}

	requestInfos = append(requestInfos, *requestInfo)

	updatedData, err := json.MarshalIndent(requestInfos, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, updatedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getBrowser(userAgent string) string {
	browserList := []string{"Chrome", "Firefox", "Safari", "Opera", "Edge", "MSIE", "Trident"}

	for _, browser := range browserList {
		if strings.Contains(userAgent, browser) {
			return browser
		}
	}

	return "Unknown"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./profiler [profiler.config]")
		os.Exit(1)
	}

	config, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	startProfilerServer(config)
}
func startProfilerServer(config Config) {
	targetURL, err := url.Parse(config.TargetServer)
	if err != nil {
		log.Fatalf("Failed to parse target server URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		headers := make(map[string]string)
		for k, v := range r.Header {
			headers[k] = strings.Join(v, ",")
		}

		ip := r.RemoteAddr
		browser := getBrowser(r.Header.Get("User-Agent"))
		received := time.Now().Format("01-02-06 - 15:04")

		requestInfo := RequestInfo{
			Method:   r.Method,
			URL:      r.URL.String(),
			Header:   headers,
			IP:       ip,
			Browser:  browser,
			Received: received,
		}

		err := writeRequestInfoToFile(&requestInfo, config.OutFile)
		if err != nil {
			log.Printf("Failed to write request info to file: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Request info saved to %s", config.OutFile)

		proxy.ServeHTTP(w, r)

	})

	address := fmt.Sprintf(":%d", config.Port)
	log.Printf("Listening on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))

}
