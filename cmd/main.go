// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/ivanmartinez/strwiki"
	"github.com/ivanmartinez/strwiki/database"
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

	// Open the database
	db := database.Connect(ctx, dbURI)
	defer db.Disconnect(ctx)

	// Start this server
	strwiki.StartServer(ctx, *url)
}
