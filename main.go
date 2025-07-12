package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	allservices "github.com/sukvij/inshorts/all-services"
	"github.com/sukvij/inshorts/inshortfers/database"
	"github.com/sukvij/inshorts/inshortfers/logs"
	"github.com/sukvij/inshorts/inshortfers/tracing"
)

// Define structs to match the expected JSON input/output
type QueryRequest struct {
	Query string `json:"query"`
}

type PredictResponse struct {
	Entities []struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"entities"`
	Concepts []string `json:"concepts"`
	Intent   string   `json:"intent"`
}

// const mlServiceURL = "http://127.0.0.1:5000/predict" // URL of your Python ML service

// func fetchResults(c *gin.Context) {
// 	var req QueryRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Prepare request for ML service
// 	jsonReqBody, err := json.Marshal(req)
// 	if err != nil {
// 		log.Printf("Error marshaling request to ML service: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare ML service request"})
// 		return
// 	}

// 	// Make POST request to ML service
// 	resp, err := http.Post(mlServiceURL, "application/json", bytes.NewBuffer(jsonReqBody))
// 	if err != nil {
// 		log.Printf("Error connecting to ML service: %v", err)
// 		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "ML service unavailable or error connecting"})
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// Read response from ML service
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("Error reading response from ML service: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response from ML service"})
// 		return
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		log.Printf("ML service returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
// 		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("ML service error: %s", string(body))})
// 		return
// 	}

// 	var mlResponse PredictResponse
// 	if err := json.Unmarshal(body, &mlResponse); err != nil {
// 		log.Printf("Error unmarshaling ML service response: %v, body: %s", err, string(body))
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ML service response"})
// 		return
// 	}

// 	// Return the ML service's response to the client
// 	c.JSON(http.StatusOK, mlResponse)
// }

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit // wit till interrupt
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
		time.Sleep(1 * time.Second)
		os.Exit(0) // Exit the program
	}()

	db, dbConnError := database.Connection()
	if dbConnError != nil {
		fmt.Println("problem with database connections... in engine file")
		return
	}
	logsForError := logs.NewAgreeGateLogger()

	tracker := tracing.InitTracer()
	fmt.Println(db, logsForError, tracker)
	app := gin.Default()
	allservices.RouteService(app, db, logsForError, tracker)
	app.Run(":8080")
}
