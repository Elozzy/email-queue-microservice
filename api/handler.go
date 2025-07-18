package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"Email-Queue-microservice/queue"
)

type Queue interface {
	Enqueue(queue.EmailJob) error
}

func HandleEmailRequest(q Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var job queue.EmailJob
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validateJob(job); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if err := q.Enqueue(job); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Job enqueued"))
	}
}

func HandleDLQ(q *queue.EmailQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dlq := q.GetDLQ()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dlq)
	}
}


func validateJob(job queue.EmailJob) error {
	if job.To == "" || job.Subject == "" || job.Body == "" {
		return errors.New("All fields are required")
	}

	emailRegex := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	if !emailRegex.MatchString(job.To) {
		return errors.New("Invalid email format")
	}

	return nil
}
