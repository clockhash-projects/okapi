# Stage 1: Build the Frontend
FROM node:24-slim AS frontend-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /app/dashboard
ARG BUILD_ID=unknown
COPY dashboard/package.json dashboard/pnpm-lock.yaml ./
RUN pnpm install --no-frozen-lockfile
COPY dashboard/ ./
RUN pnpm build

# Stage 2: Build the Go Backend
FROM library/golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
# Copy only necessary source files to maintain layer cache
COPY main.go ./
COPY adapters/ ./adapters/
COPY api/ ./api/
COPY internal/ ./internal/
# Copy the built frontend into the dashboard/dist directory for embedding
COPY --from=frontend-builder /app/dashboard/dist ./dashboard/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o okapi main.go

# Stage 3: Final Minimal Image
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
# Create a non-root user for security
RUN addgroup -S okapi && adduser -S okapi -G okapi
USER okapi

COPY --from=backend-builder /app/okapi .
# Copy default config, adapters, and dashboard
COPY config.example.yaml ./config.yaml
COPY adapters/config ./adapters/config
COPY --from=backend-builder /app/dashboard/dist ./dashboard/dist

EXPOSE 8080
CMD ["./okapi"]
