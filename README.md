# 🛡️ API Guardian (Zero Trust Gateway)

![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg) ![Redis](https://img.shields.io/badge/Redis-Stack-red.svg) ![Docker](https://img.shields.io/badge/Docker-Compose-2496ED.svg) ![License](https://img.shields.io/badge/License-MIT-green.svg)

**API Guardian** is a high-performance API Gateway & Reverse Proxy solution built using **Go (Golang)** and **Redis**.

This project is designed as the **First Line of Defense** to protect backend services from unauthorized access, DDoS attacks, and traffic abuse. Adopting **Zero Trust** principles, the system verifies every single incoming request before forwarding it to the main server.

---

## 🚀 Key Features

### 🔒 Security Core

- ✅ **Zero Trust Architecture**: No request is trusted by default; everything must be validated.
- ✅ **API Key Validation**: The gate opens only for clients with valid credentials (`X-API-KEY`).
- ✅ **Distributed Token Bucket Algorithm**: Implemented via Redis Lua Scripting to ensure atomic operations. This prevents "race conditions" in high-traffic distributed environments.
- ✅ **State Persistence**: Rate limit counters and buckets are stored in Redis, meaning the system state survives application restarts (No "amnesia" on reboot).
- ✅ **Secure Reverse Proxy**: Hides the original server identity and manages HTTP headers automatically.
- ✅ **Fail-Closed Security Logic**: Designed with a "Security First" mindset. If the Redis connection drops, the system defaults to a safe state to protect backend integrity.
- ✅ **Violation Counter**: Tracks repeated limit breaches (429 Too Many Requests) using a dedicated violation:{ip} registry in Redis.
- ✅ **Dynamic Blacklist**: Automatically escalates repeat offenders to a "Blacklist" state once they exceed the MAX_VIOLATIONS threshold.
- ✅ **Basic Web Application Firewall (WAF)**: Inspecting every incoming HTTP request in real-time before it can reach the backend services or the database.
- ✅ **IP Whitelisting**: A special lane (VVIP) for administrator IPs or internal services to bypass limitations.
- ✅ **Circuit Breakers**: Adds system intelligence by automatically cutting off traffic to a failing backend.
- ✅ **PII Redaction/Masking**: Automatically turns sensitive data like Emails or IDs into `***`.

### 📊 Observability & Operations

- ✅ **GeoIP & Device Intelligence**: Automatically detects User Country, City, OS, Browser, and Bot status for forensic analysis.
- ✅ **Smart Response Caching (Redis)**: intercepts repeated requests for the same data and serves them directly from memory, bypassing the expensive process of re-querying the database or re-processing business logic.
- ✅ **Relational Audit Logging**: Every security event is synchronized to a PostgreSQL database, enabling long-term trend analysis and legal compliance.
- ✅ **Prometheus Metrics Exporter**: Exposes internal system health and performance data through a `/metrics` endpoint, allowing DevOps teams to scrape, store, and visualize real-time data using the industry-standard Prometheus & Grafana stack.
- ✅ **Structured Audit Logging**: Uses Zerolog for high-performance, machine-readable logs, making the system 100% ready for ELK Stack or Grafana Loki integration.
- ✅ **Log Rotation**: Automatic log file management (Lumberjack) to prevent disk saturation.
- ✅ **Round-Robin Load Balancer**: Enables the gateway to manage a pool of multiple backend instances, ensuring the system remains operational even if one server fails.
- ✅ **Graceful Shutdown**: Handles `SIGTERM` signals to terminate services without cutting off active connections.
- ✅ **Dockerized**: Ready to run in any container environment with a single command.

---

## 🏗️ System Architecture

Here is the workflow of how API Guardian protects your Backend:

1.  **Client** sends a request to API Guardian (Port 8080).
2.  **Guardian** checks key validity and quota limits in **Redis** (In-Memory).
3.  **Denied**: If invalid/limit exceeded, Guardian immediately replies with a JSON Error (401/429).
4.  **Allowed**: If safe, the request is forwarded to the `TARGET_URL`.
5.  **Audit**: Every transaction is recorded in `logs/app.log`.

---

## 📂 Project Structure

```text
api-guardian/
├── cmd/
│   └── server/             # Application entry point (main.go)
├── internal/               # Private application and library code
│   ├── middleware/         # WAF, Caching, PII Masking, Rate Limiter logic
│   └── usecase/            # Business logic & Database interfaces
├── configs/                # Configuration files
│   └── geoip/              # MaxMind GeoIP Database (.mmdb)
├── deployments/            # Docker, Compose & Deployment configs
├── logs/                   # Audit logs (Automatically created)
├── web/                    # (In Progress) React Dashboard UI
├── .env.example            # Environment template
└── go.mod                  # Dependency management
```

## 🛠️ Installation & Setup

### ⚠️ Pre-requisites (GeoIP)

This project requires the `MaxMind GeoLite2` database for intelligence features and `Docker` to run (remommended).

- Download `GeoLite2-City.mmdb` from `MaxMind`.
- Place the file inside: `configs/geoip/`
  _If this file is missing, location data will show as "Unknown"._

#### 1. Clone & Prep

```bash
git clone https://github.com/niba17/api-guardian.git
cd api-guardian
go mod vendor
```

#### 2. Environment

```bash
cp configs/.env.example .env
```

_remember to set your target url `TARGET_URL` on `.env`, also you can provide multiple target url using a comma (,) without spaces. Example: `TARGET_URL=url1,url2,url3`._

#### 3. Run (Docker Compose)

```bash
docker-compose -f deployments/docker-compose.yml up -d --build
```

---

#### 4. Testing

```bash
curl.exe -i http://localhost:8080/status
```

_success response example:_

```text
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 11 Feb 2026 15:01:43 GMT
Content-Length: 118

{"circuit_breaker":"Closed (Normal)","redis_connection":"Connected","system":"Healthy","time":"2026-02-11T15:01:43Z"}
```
