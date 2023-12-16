package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fredrikaverpil/promethea/internal/servers"
)

func main() {
	// OLAMA_URL environment variable
	ollamaUrl := os.Getenv("OLLAMA_URL")
	if ollamaUrl == "" {
		ollamaUrl = "http://127.0.0.1:11434"
		log.Printf("OLLAMA_URL not set, defaulting to %s", ollamaUrl)
	}

	// Create channel for graceful shutdown
	stopChan := make(chan os.Signal, 1)

	restServer := servers.NewRESTServer(":8080", ollamaUrl)
	go func() {
		if err := restServer.Start(); err != nil {
			log.Fatalf("%v", err)
		}
	}()

	// Handle graceful shutdown
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-stopChan
		restServer.Stop()
		os.Exit(0)
	}()

	// Wait forever
	select {}
}
