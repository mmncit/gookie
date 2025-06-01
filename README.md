# 🟢 Simple Go HTTP Server

A minimal HTTP server written in Go, using the standard [`net/http`](https://pkg.go.dev/net/http) package. It includes a basic health check endpoint and sample handlers to demonstrate routing.

---

## ✅ Features

- Lightweight
- Built with idiomatic Go (Go 1.18+)
- Includes a `/health` endpoint
- Simple GET and POST examples

---

## 🚀 Prerequisites

- Go **1.18 or higher** installed

---

## 🛠️ Running the Server

By default, the server listens on **port 8080**.

To run the server:

```bash
go run ./main.go -port 8080
```

---

## 📡 Example Endpoints

### 🔍 Health Check

```bash
curl http://localhost:8080/health
```

### 📥 POST Example

```bash
curl -X POST http://localhost:8080/events
```

### 📤 GET Example

```bash
curl http://localhost:8080/events/2
```

---

## 📄 License

This project is licensed under the [MIT License](./LICENSE).
