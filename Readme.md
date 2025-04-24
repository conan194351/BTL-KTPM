# 🛠️ Project Setup

This project provides a simple setup for local development using Docker and Make.

---

## 🚀 Getting Started

### Step 1: Install Go

Ensure you have **Go version 1.23.5** installed.

You can download it from the official site:  
👉 [https://go.dev/dl/go1.23.5](https://go.dev/dl/go1.23.5)

To verify your installation:

```bash
go version
# Expected: go version go1.23.5 linux/amd64
```

### Step 2: Start Services with Docker

```bash
cd deployment
docker-compose up -d
```

### Step 3: Run the Application
```bash
make local
```

## Project Structure

```plaintext
.
├── cmd/                # Main application entry points
├── internal/           # Application core logic
├── deployment/         # Docker and docker-compose setup
│   └── docker-compose.yml
├── pkg/                # Shared libraries and packages
├── Makefile            # Make commands for local dev
└── README.md           # Project setup and usage guide
```