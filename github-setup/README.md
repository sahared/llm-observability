# 🚀 LLM Observability Platform

> Production-ready observability, monitoring, and debugging platform for LLM-powered applications and AI agents.

[![Backend CI](https://github.com/YOUR_USERNAME/llm-observability/workflows/Backend%20CI/badge.svg)](https://github.com/YOUR_USERNAME/llm-observability/actions)
[![Frontend CI](https://github.com/YOUR_USERNAME/llm-observability/workflows/Frontend%20CI/badge.svg)](https://github.com/YOUR_USERNAME/llm-observability/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev)
[![Node Version](https://img.shields.io/badge/Node-18+-339933?logo=node.js)](https://nodejs.org)

---

## 📋 Overview

**LLM Observability Platform** is a comprehensive "DataDog for AI Agents" - built to solve the unique challenges of monitoring non-deterministic LLM systems. Track costs, performance, and reliability across OpenAI, Anthropic, Cohere, and custom models.

### 🎯 The Problem

- **80% of AI projects fail** due to infrastructure immaturity
- **LLM costs are unpredictable** and can spiral out of control
- **Debugging AI systems** requires specialized tools
- **No visibility** into multi-step agent workflows

### ✨ Our Solution

- 🔍 **Debug Faster** - Full distributed tracing with context preservation
- 💰 **Optimize Costs** - Track spending by team, project, and model
- ⚠️ **Prevent Failures** - Real-time alerting for anomalies
- 📊 **Understand Behavior** - Analytics across all prompts and responses
- 🎯 **Improve Quality** - Track performance and reliability

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────┐
│         Client Applications                 │
│    (Python/JS/Go SDKs in your code)        │
└──────────────────┬──────────────────────────┘
                   │ HTTPS
                   ▼
┌─────────────────────────────────────────────┐
│        Backend API (Go + Fiber)             │
│   Ingest • Query • Analytics • Auth        │
└──────────────────┬──────────────────────────┘
                   │
    ┌──────────────┼──────────────┐
    │              │              │
    ▼              ▼              ▼
┌────────┐   ┌─────────┐   ┌──────────┐
│ClickH. │   │  Redis  │   │  Kafka   │
│(Store) │   │ (Cache) │   │ (Queue)  │
└────────┘   └─────────┘   └──────────┘
                               │
                               ▼
                         ┌──────────┐
                         │Prometheus│
                         │ Grafana  │
                         └──────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- **Docker Desktop** v24.0+
- **Go** v1.21+
- **Node.js** v18+
- **Python** v3.11+ (for SDK)

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/YOUR_USERNAME/llm-observability.git
cd llm-observability

# 2. Start infrastructure services
cd infrastructure
docker-compose up -d

# 3. Initialize database
docker exec -i llm-obs-clickhouse clickhouse-client --multiquery < ../backend/migrations/001_initial_schema.up.sql

# 4. Start backend
cd ../backend
cp .env.example .env
go run cmd/api/main.go

# 5. Start frontend (new terminal)
cd ../frontend
npm install
npm run dev
```

### Access the Platform

- **Frontend Dashboard**: http://localhost:5173
- **Backend API**: http://localhost:8080/health
- **Grafana**: http://localhost:3001 (admin/admin)
- **Prometheus**: http://localhost:9090

### Send Your First Trace

```bash
curl -X POST http://localhost:8080/api/v1/traces \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo-key" \
  -d '{
    "organization_id": "org-demo",
    "project_id": "proj-demo",
    "model": "gpt-4",
    "provider": "openai",
    "trace_type": "single_call",
    "spans": [{
      "name": "chat_completion",
      "model": "gpt-4",
      "provider": "openai",
      "input": "What is observability?",
      "output": "Observability is the ability to measure...",
      "prompt_tokens": 10,
      "completion_tokens": 25,
      "duration_ms": 850,
      "status": "success"
    }]
  }'
```

---

## ✨ Features

### 🔬 Core Capabilities

- **Distributed Tracing** - Track multi-step LLM workflows
- **Cost Attribution** - Per-model, per-team cost tracking
- **Performance Monitoring** - P50, P95, P99 latency metrics
- **Real-time Analytics** - Live dashboards and insights
- **Anomaly Detection** - ML-powered alerting
- **Semantic Search** - Search across prompts and responses

### 🤖 Supported Providers

- ✅ **OpenAI** - GPT-4, GPT-4 Turbo, GPT-3.5
- ✅ **Anthropic** - Claude 3.5 Sonnet, Claude 3 Opus/Sonnet/Haiku
- ✅ **Cohere** - Command, Embed, Rerank
- 🔄 **Google** - Gemini (coming soon)
- ✅ **Custom** - Self-hosted models

---

## 🛠️ Technology Stack

### Backend
- **Go 1.21** + **Fiber** - High-performance HTTP server
- **ClickHouse** - Time-series analytics database
- **Apache Kafka** - Event streaming
- **Redis** - Caching layer

### Frontend
- **React 18** + **TypeScript** - Modern UI framework
- **Tailwind CSS** - Utility-first styling
- **Recharts** - Data visualization
- **Zustand** - State management

### Infrastructure
- **Docker** + **Docker Compose** - Container orchestration
- **Prometheus** + **Grafana** - Monitoring
- **GitHub Actions** - CI/CD

---

## 📁 Project Structure

```
llm-observability/
├── backend/               # Go backend services
│   ├── cmd/api/          # Application entry point
│   ├── internal/         # Private application code
│   │   ├── api/         # HTTP handlers
│   │   ├── models/      # Data models
│   │   ├── services/    # Business logic
│   │   ├── repository/  # Data access layer
│   │   └── middleware/  # HTTP middleware
│   └── migrations/       # Database migrations
├── frontend/             # React frontend
│   ├── src/
│   │   ├── pages/       # Page components
│   │   ├── components/  # Reusable components
│   │   ├── services/    # API clients
│   │   └── stores/      # State management
│   └── public/          # Static assets
├── sdk/                  # Python SDK (coming soon)
├── infrastructure/       # Infrastructure as Code
│   ├── docker-compose.yml
│   └── clickhouse/      # DB configs
└── docs/                # Documentation
```

---

## 🧪 Testing

```bash
# Backend tests
cd backend
go test -v ./...

# Frontend tests
cd frontend
npm test

# SDK tests (coming soon)
cd sdk
pytest
```

---

## 📊 Development Status

### ✅ Phase 1: Foundation (Current)
- [x] Docker Compose infrastructure
- [x] ClickHouse schema
- [x] Backend API structure
- [x] Frontend scaffolding
- [x] CI/CD workflows
- [ ] Real database operations
- [ ] Functional trace ingestion
- [ ] Dashboard with live data

### 🚧 Phase 2: Core Features (Next)
- [ ] Complete dashboard UI
- [ ] Trace list and detail pages
- [ ] Real-time updates
- [ ] Authentication system

### 📋 Upcoming Phases
- Phase 3: Python SDK
- Phase 4: Advanced Analytics
- Phase 5: Cloud Deployment
- Phase 6: Monitoring & Alerting
- Phase 7: Testing & Optimization
- Phase 8: Documentation & Launch

---

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

---

## 📝 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- Inspired by [OpenTelemetry](https://opentelemetry.io/), [Jaeger](https://www.jaegertracing.io/), and [DataDog](https://www.datadoghq.com/)
- Built for the AI/ML community

---

## 📞 Contact

- **GitHub Issues**: [Report bugs or request features](https://github.com/YOUR_USERNAME/llm-observability/issues)
- **Discussions**: [Join the conversation](https://github.com/YOUR_USERNAME/llm-observability/discussions)

---

## 🗺️ Roadmap

### v1.0 (Phase 1-2) - Foundation
- Core tracing functionality
- Real-time dashboard
- Basic analytics

### v1.1 (Phase 3-4) - SDK & Analytics
- Python SDK
- Advanced cost optimization
- A/B testing framework

### v1.2 (Phase 5-6) - Production Ready
- Cloud deployment
- Kubernetes support
- Advanced monitoring

### v2.0 (Future)
- JavaScript/TypeScript SDK
- Mobile dashboard
- AI-powered insights

---

**⭐ Star this repository if you find it useful!**

Made with ❤️ by [Your Name](https://github.com/YOUR_USERNAME)
