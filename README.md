# AgentTrace

**AgentTrace** is a modular backend system for logging, visualizing, and evaluating LLM-based agent behavior. It allows tracing of agent execution, capturing prompts, outputs, latency, token usage, and substeps — forming the foundation for deeper observability and evaluation.

---

## ✨ Features

- 🚀 Modular architecture (Gin + MongoDB)
- 📥 POST endpoint to ingest agent traces
- 📦 Clean separation: handlers, repositories, routers, config, and assembler
- 🔧 Dockerized development with MongoDB
- ✅ Ready for unit testing with interfaces and mocks

---

## 🧱 Tech Stack

- **Go 1.24+**
- **Gin** (HTTP framework)
- **MongoDB** (trace persistence)
- **Docker + Compose** (local setup)
- **Testify** (mocking & assertions)
- **envconfig** (env vars management)

---

## 📦 Project Structure
```
agent-trace/
├── assembler/    # Assembles app dependencies
├── cmd/          # App entrypoint (main.go)
├── config/       # Env config loading via envconfig
├── internal/
│   ├── db/       # MongoDB client init
│   ├── handler/  # HTTP handlers (interface + implementation)
│   ├── model/    # Domain models (Trace, Substep, etc.)
│   ├── repository/ # Mongo repo + interface
│   └── router/   # Route setup and separation
├── test/         # Unit tests (e.g., handler with mocks)
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🚀 Getting Started

### 1. Clone the repo
```bash
git clone git@github.com:zkropotkine/agent-trace.git
cd agent-trace
```

### 2. Build and run with Docker Compose
```bash
docker-compose up --build
```

This will:
* Start MongoDB on `localhost:27017`
* Start the app on `localhost:8080`

## 📬 Example API Usage

### `POST /api/traces`

Ingest a trace:
```bash
curl -X POST http://localhost:8080/api/traces \
  -H "Content-Type: application/json" \
  -d '{
    "trace_id": "abc123",
    "session_id": "session-xyz",
    "agent_name": "DocumentAgent",
    "status": "success",
    "input_prompt": "Summarize privacy policy.",
    "output": "Privacy policy summary...",
    "latency_ms": 350,
    "token_usage": {
      "input_tokens": 100,
      "output_tokens": 200,
      "total": 300
    },
    "substeps": [
      {
        "name": "Retriever",
        "input": "privacy",
        "output": "[doc1, doc2]",
        "status": "success",
        "start": "2025-05-01T01:23:00Z",
        "end": "2025-05-01T01:23:01Z"
      }
    ]
  }'
```

### ✅ Expected response
```json
{
  "message": "trace saved"
}
```

## ⚙️ Configuration

| Env Variable | Default | Description |
|--------------|---------|-------------|
| `AGENTTRACE_PORT` | `8080` | Port for the HTTP server |
| `AGENTTRACE_MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `AGENTTRACE_ENV` | `development` | App environment |

Set these in your shell or use `.env` + tools like `direnv`.

## 🧪 Running Tests
```bash
go test ./...
```

## 📌 Roadmap

* POST /api/traces
* GET /api/traces/:id
* Web dashboard for trace visualization
* LLM-based evaluation engine (AgentManager)
* Auth and team-based trace access
* Deployment pipeline

## 🧑‍💻 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss the design or scope.

## 📄 License

MIT

## 💡 Inspiration

AgentTrace is inspired by the emerging need for LLM agent observability, evaluation, and runtime introspection in autonomous or tool-using systems.

---

Let me know if you'd like:
- A badge/header setup (build passing, version, etc.)
- A makefile or task runner
- A starter .env or .env.example file for developer setup