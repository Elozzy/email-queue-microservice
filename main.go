package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"Email-Queue-microservice/api"
	"Email-Queue-microservice/queue"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := queue.Config{
		WorkerCount: 3,
		QueueSize:   10,
	}

	emailQueue := queue.NewEmailQueue(config)
	emailQueue.StartWorkers(ctx)

	http.HandleFunc("/send-email", api.HandleEmailRequest(emailQueue))
	http.HandleFunc("/dlq", api.HandleDLQ(emailQueue))


	server := &http.Server{Addr: ":8080"}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		log.Println("\nReceived shutdown signal")
		server.Shutdown(ctx)
		cancel()
		emailQueue.Close()
	}()

	log.Println("🚀 Email Queue Service started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	emailQueue.Wait()
	log.Println("✅ All workers finished. Exiting.")
}
