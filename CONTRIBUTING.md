# Contributing to Okapi

Thank you for your interest in contributing! Okapi thrives on community-added service adapters.

## Adding a New Service Adapter (YAML)

Most status pages (GitHub, Stripe, etc.) use Atlassian Statuspage. These can be added without writing any Go code.

1.  **Fork the repository.**
2.  **Identify the status subdomain:** Visit the service's status page and check if it has a `/api/v2/summary.json` endpoint.
3.  **Create a new YAML file** in `adapters/config/<service-id>.yaml`:
    ```yaml
    id: my-service
    display_name: My Service Name
    kind: statuspage
    subdomain: status.myservice.com
    poll_interval_seconds: 60
    ```
4.  **Test your adapter:**
    ```bash
    go run main.go
    # Verify via curl http://localhost:8080/health/my-service
    ```
5.  **Submit a Pull Request.**

## Adding a Go Adapter

For services with custom APIs (like AWS or Azure):

1.  Create a new file in `adapters/code/<service-id>.go`.
2.  Implement the `HealthAdapter` interface.
3.  Register it in `main.go`.
4.  Include unit tests in `adapters/code/<service-id>_test.go`.

## Development Environment

- **Hot Reload:** Use `air`.
- **Testing:** Run `go test ./...`.
- **Linting:** We follow standard Go formatting (`go fmt`).

## Code of Conduct

Please be respectful and professional in all interactions within this project.
