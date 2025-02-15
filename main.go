package main

import (
	"crypto/tls"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	// Interval for updating service status
	updateInterval = 60 * time.Second
)

// Service represents a monitored service
type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
	IsUp bool   `json:"isUp"`
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	IsDown    bool      `json:"isDown"`
	StartTime time.Time `json:"startTime"`
	Name      string    `json:"name"`
}

// Config represents the configuration file structure
type Config struct {
	Services []Service `json:"services"`
}

var (
	// Map to store service downtimes
	serviceDowntimes = make(map[string]*ServiceStatus)
	// Mutex for synchronizing access to serviceDowntimes
	mutex sync.RWMutex
	// Config object to store the loaded configuration
	config Config
)

func main() {
	// Load the configuration file
	if err := loadConfig(); err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Start the service checker
	startServiceChecker()

	// Serve static files from the "icons" directory
	fs := http.FileServer(http.Dir("icons"))
	http.Handle("/icons/", http.StripPrefix("/icons/", fs))
	// Serve the CSS file
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "style.css")
	})
	// Handle the status endpoint
	http.HandleFunc("/status", getServiceStatus)
	// Handle the root endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("template.html"))
		tmpl.Execute(w, config)
	})

	// Start the HTTP server
	log.Println("Is It Dead Yet? running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// startServiceChecker starts the service status checker
func startServiceChecker() {
	// Initial check of all services
	checkAllServices()

	// Set up a ticker to check services periodically
	ticker := time.NewTicker(updateInterval)
	go func() {
		for range ticker.C {
			checkAllServices()
		}
	}()
}

// checkServiceStatus checks the status of a service
func checkServiceStatus(url string) (time.Time, bool) {
	client := http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
			IdleConnTimeout:   1 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return time.Time{}, false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return time.Time{}, false
	}
	if resp != nil {
		defer resp.Body.Close()
		return time.Time{}, resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	return time.Time{}, false
}

// checkAllServices checks the status of all services
func checkAllServices() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, service := range config.Services {
		_, isUp := checkServiceStatus(service.URL)
		if !isUp {
			if _, exists := serviceDowntimes[service.URL]; !exists {
				serviceDowntimes[service.URL] = &ServiceStatus{
					IsDown:    true,
					StartTime: time.Now(),
					Name:      service.Name,
				}
			}
		} else {
			delete(serviceDowntimes, service.URL)
		}
	}
}

// getServiceStatus handles the /status endpoint
func getServiceStatus(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	response := struct {
		Services  []Service                 `json:"services"`
		Downtimes map[string]*ServiceStatus `json:"downtimes"`
		Timestamp int64                     `json:"timestamp"`
	}{
		Services:  make([]Service, len(config.Services)),
		Downtimes: serviceDowntimes,
		Timestamp: time.Now().Unix(),
	}

	for i, service := range config.Services {
		isUp := true
		if _, exists := serviceDowntimes[service.URL]; exists {
			isUp = false
		}
		response.Services[i] = Service{
			Name: service.Name,
			URL:  service.URL,
			Icon: service.Icon,
			IsUp: isUp,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// loadConfig loads the configuration file
func loadConfig() error {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &config)
}
