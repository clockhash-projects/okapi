# Okapi

**The Unified Health Check Proxy & API for 80+ Global Services.**

Okapi is a high-performance, open-source proxy service that provides a single, centralized API for querying the health status of over 80 third-party providers, including AWS, GitHub, Stripe, and Cloudflare.

Instead of integrating with dozens of different status page schemas and polling intervals, Okapi offers a unified, normalized interface. It handles the complexities of fetching, parsing, and caching health data, allowing your systems to monitor all their dependencies through a single endpoint.

This is the open-source version of Okapi. For internal deployments, please refer to the `okapi-internal` repository which may contain additional configurations or scripts.

---

## 🚀 Key Features

- **Unified API:** One consistent schema for 80+ providers (and counting).
- **Federated Status:** Query status for multiple services in a single batch request.
- **Pluggable Architecture:** Easily add new services via YAML config or custom Go code.
- **Embedded Dashboard:** A built-in React-based dashboard for real-time visual monitoring.
- **Webhooks:** Register callback URLs to receive instant notifications when a service status changes.
- **Intelligent Caching:** Support for both Redis and In-Memory caching to reduce upstream latency and avoid rate limits.
- **Production Ready:** Includes API authentication, structured logging (Zap), and Prometheus metrics.

---

## 🏗️ Architecture

Okapi is built on a modular "Adapter" system:

1.  **Go Adapters (`adapters/code/`):** For services with complex APIs or non-standard status pages (e.g., AWS, GCP, Azure, Slack).
2.  **YAML Adapters (`adapters/config/`):** For 80+ services that use standard Statuspage.io, RSS, or HTTP-based status reporting.

A background polling worker periodically refreshes the status of all registered services, updates the cache, and triggers webhooks if any status changes are detected.

---

## 🚦 Getting Started

### 1. Installation
Ensure you have Go 1.24+ installed.

To get started with the open-source version:
```bash
git clone <your-public-repo-url> okapi
cd okapi
make build
```
Replace `<your-public-repo-url>` with the actual URL of the public repository once it's hosted.

### 2. Configuration
Copy the example configuration and adjust it to your needs.

```bash
cp config.example.yaml config.yaml
```
Environment variables can be used to override configuration settings (e.g., `OKAPI_PORT=9090`).

### 3. Run
```bash
./okapi
```
The API will be available at `http://localhost:8080/api` and the dashboard at `http://localhost:8080/`.

---

## 🛠️ API Reference

All API endpoints are prefixed with `/api`.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/_health` | Self-health check for Okapi |
| **GET** | `/services` | List all 80+ registered service IDs |
| **GET** | `/health/{service}` | Get the current status for a single service (e.g., `/health/github`) |
| **GET** | `/health` | Batch status check (e.g., `?services=aws,github,stripe`) |
| **GET** | `/incidents` | Consolidated view of all services currently experiencing issues |
| **GET** | `/maintenance` | Consolidated view of all upcoming scheduled maintenance |
| **GET** | `/webhooks` | List all active webhook subscriptions |
| **POST** | `/webhooks` | Register a new webhook for status change notifications |
| **DELETE**| `/webhooks/{id}` | Remove a registered webhook |
| **GET** | `/metrics` | Prometheus metrics for monitoring Okapi's health |
| **GET** | `/help` | Detailed API documentation and endpoint list |

---

## 🔌 Supported Adapters (80+)

Okapi supports a wide range of services across various categories:

- **Cloud Infrastructure:** AWS, GCP, Azure, Cloudflare, DigitalOcean, Heroku, Vercel, Netlify.
- **CI/CD & Dev Tools:** GitHub, GitLab, CircleCI, Travis CI, Bitbucket, Docker, npm, PyPI.
- **Data & Databases:** MongoDB, Redis, Supabase, PlanetScale, Snowflake, Aiven, InfluxDB.
- **Communication:** Slack, Discord, Twilio, Zoom, SendGrid, Mailgun, Pusher.
- **SaaS & Platforms:** Stripe, Shopify, OpenAI, Anthropic, Airtable, HubSpot, Zendesk, Intercom.
- **Observability:** Datadog, Grafana, New Relic, Sentry, PagerDuty, Better Stack.

*See `adapters/config/` for the full list of declarative adapters.*

---

## 🔗 Webhooks

Okapi can notify your systems whenever a service changes its status. To register a webhook:

```bash
POST /api/webhooks
{
  "url": "https://your-app.com/webhooks/okapi",
  "services": ["aws", "github"],
  "secret": "your-webhook-secret"
}
```

---

## 📊 Dashboard

Okapi includes a modern React-based dashboard that is embedded directly into the Go binary. It provides:
- Real-time status overview of all registered services.
- Detailed component-level health for complex providers.
- Incident history and scheduled maintenance timelines.

Access it by navigating to the root URL of your Okapi instance.

---

## 🛠️ Development

### Adding a New Service
1. **Simple Statuspage/RSS:** Add a `.yaml` file to the appropriate subdirectory in `adapters/config/`.
2. **Complex API:** Create a new adapter in `adapters/code/`, implement the `HealthAdapter` interface, and register it in `main.go`.

### Running Tests
```bash
make test         # Run all unit tests
make e2e          # Run end-to-end tests (requires Docker)
make lint         # Run golangci-lint
```

---

## ⚖️ License

Okapi is distributed under the **MIT License**. See `LICENSE` for more information.
