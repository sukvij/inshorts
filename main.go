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

func main() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit // wait till interrupt
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
