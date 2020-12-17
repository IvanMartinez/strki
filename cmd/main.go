package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ivanmartinez/boocat/database"
	"github.com/ivanmartinez/boocat/formats"
	"github.com/ivanmartinez/boocat/server"
)

func main() {
	// Parse flags
	url := flag.String("url", "localhost:80", "This server's base URL")
	dbURI := flag.String("dburi", "mongodb://127.0.0.1:27017", "Database URI")
	flag.Parse()

	// Create channel for listening to OS signals and connect OS interrupts to
	// the channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("received signal %v", oscall)
		cancel()
	}()

	// Start services
	db := database.Connect(ctx, dbURI, []string{"author", "book"})
	formats.Initialize(db)
	server.Initialize(ctx, db, *url)
	server.Start()

	// Wait for ctx to be cancelled
	<-ctx.Done()

	// New context with timeout to shut the HTTP server down
	ctxShutDown, cancel := context.WithTimeout(context.Background(),
		5*time.Second)

	// Shut services down
	server.ShutdownServer(ctxShutDown)
	db.Disconnect(ctxShutDown)
}
