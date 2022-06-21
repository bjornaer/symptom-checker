package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bjornaer/sympton-checker/ent"
	_ "github.com/mattn/go-sqlite3"
)

// returns dbAddr, dbBackend, portsFileName
func getEnvVars() string {
	symptomsFile := os.Getenv("SYMPTOMS_FILE")
	if len(symptomsFile) == 0 {
		symptomsFile = "http://www.orphadata.org/data/xml/en_product4.xml"
	}
	return symptomsFile
}

func main() {
	// create a new ent client with an in memory database.
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Panicf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	symptomsFile := getEnvVars()
	// dbClient, err := db.InitDBClientWithRetry(5, dbBackend, dbAddr)
	// if err != nil {
	// 	log.Fatal("Failed initialising DB connection")
	// }
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Panicf("failed creating schema resources: %v", err)
	}
	s := NewServer(client)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: s,
	}

	s.LoadSymptomsFromRemote(symptomsFile)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
