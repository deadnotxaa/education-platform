package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL          string
	httpClient       *http.Client
	usedReportLimits []int
	mu               sync.Mutex
	rng              *rand.Rand
}

type ReportRequest struct {
	LimitNumber int `json:"limit_number"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		usedReportLimits: make([]int, 0, 10),
		rng:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) logRequest(endpoint, params string, responseTime time.Duration, statusCode int) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] %s | Params: %s | Time: %.2fms | Status: %d\n",
		timestamp, endpoint, params, float64(responseTime.Nanoseconds())/1e6, statusCode)
}

func (c *Client) getCourse() {
	courseID := c.rng.Intn(1000) + 1
	url := fmt.Sprintf("%s/v1/course/getcourse", c.baseURL)

	// Create JSON body as expected by the backend
	requestData := map[string]int{
		"id": courseID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Create GET request with body (unusual but that's what your backend expects)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		log.Printf("Error in getCourse: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	c.logRequest("GET /v1/course/getcourse", fmt.Sprintf("id=%d", courseID), responseTime, resp.StatusCode)
}

func (c *Client) getUser() {
	userID := c.rng.Intn(1000) + 1
	url := fmt.Sprintf("%s/v1/user/getuser", c.baseURL)

	// Create JSON body as expected by the backend
	requestData := map[string]int{
		"id": userID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Create GET request with body (unusual but that's what your backend expects)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		log.Printf("Error in getUser: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	c.logRequest("GET /v1/user/getuser", fmt.Sprintf("id=%d", userID), responseTime, resp.StatusCode)
}

func (c *Client) getTopCoursesReport() {
	c.mu.Lock()
	var limitNumber int
	var cacheIndicator string

	// 30% chance to reuse a previously used limit for cache testing
	if len(c.usedReportLimits) > 0 && c.rng.Float32() < 0.3 {
		limitNumber = c.usedReportLimits[c.rng.Intn(len(c.usedReportLimits))]
		cacheIndicator = "(cached)"
	} else {
		limitNumber = c.rng.Intn(1000) + 1
		cacheIndicator = "(fresh)"

		// Store this limit for future cache testing (keep only last 10)
		c.usedReportLimits = append(c.usedReportLimits, limitNumber)
		if len(c.usedReportLimits) > 10 {
			c.usedReportLimits = c.usedReportLimits[1:]
		}
	}
	c.mu.Unlock()

	url := fmt.Sprintf("%s/v1/report/get-top-courses-report", c.baseURL)

	requestData := ReportRequest{
		LimitNumber: limitNumber,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}

	// Create GET request with body (unusual but that's what your backend expects)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	responseTime := time.Since(start)

	if err != nil {
		log.Printf("Error in getTopCoursesReport: %v", err)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()

	params := fmt.Sprintf("limit=%d %s", limitNumber, cacheIndicator)
	c.logRequest("GET /v1/report/get-top-courses-report", params, responseTime, resp.StatusCode)
}

func (c *Client) frequentRequestsWorker() {
	for {
		// Randomly choose between course and user requests
		if c.rng.Intn(2) == 0 {
			c.getCourse()
		} else {
			c.getUser()
		}

		// Sleep for 50-200ms between frequent requests
		sleepTime := time.Duration(c.rng.Intn(150)+50) * time.Millisecond
		time.Sleep(sleepTime)
	}
}

func (c *Client) reportRequestsWorker() {
	for {
		c.getTopCoursesReport()
		// Sleep for 3-4 seconds between report requests
		sleepTime := time.Duration(c.rng.Intn(1000)+3000) * time.Millisecond
		time.Sleep(sleepTime)
	}
}

func (c *Client) Run() {
	fmt.Println("Starting Backend Load Testing Client...")
	fmt.Printf("Target URL: %s\n", c.baseURL)
	fmt.Println(strings.Repeat("=", 80))

	// Wait a bit for backend to be ready
	time.Sleep(5 * time.Second)

	// Start worker goroutines
	go c.frequentRequestsWorker()
	go c.reportRequestsWorker()

	// Keep main goroutine alive
	select {} // Block forever
}

func main() {
	baseURL := "http://backend:8080"
	client := NewClient(baseURL)
	client.Run()
}
