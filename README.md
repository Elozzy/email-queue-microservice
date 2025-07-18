# email-queue-microservice
# 📬 Go Email Queue Microservice

A lightweight Go microservice that accepts email jobs over HTTP, queues them, and processes them asynchronously using a concurrent worker pool.

---

## 🚀 Features

- ✅ HTTP API to enqueue email jobs
- ✅ In-memory job queue with channels
- ✅ Multiple concurrent workers
- ✅ Simulated email delivery
- ✅ Graceful shutdown on SIGINT/SIGTERM
- ✅ Retry failed jobs (up to 3 times)
- ✅ Dead Letter Queue (DLQ) for permanently failed jobs

---

## 📦 Folder Structure

emailservice/
├── main.go # Entry point
├── api/
│ └── handler.go # HTTP handlers
├── queue/
│ ├── config.go # Queue config
│ ├── interface.go # Job interface
│ ├── job.go # Job model
│ └── queue.go # Queue logic


---

## ⚙️ Configuration

Configuration is currently hardcoded in `main.go`:

```go
cfg := queue.Config{
    QueueSize:   10,
    WorkerCount: 3,
}

## How to Run
go run main.go

## 📬 API Endpoints

### 1. `POST /send-email`

Enqueue a new email job.

#### 📝 Request Payload

```json
{
  "to": "user@example.com",
  "subject": "Welcome!",
  "body": "Thanks for signing up."
}

curl -X POST http://localhost:8080/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "to": "user@example.com",
    "subject": "Welcome!",
    "body": "Thanks for signing up."
}'
