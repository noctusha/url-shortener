# ---------- build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY config /app/config

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o url-shortener ./cmd/url-shortener

# ---------- runtime stage ----------
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/url-shortener /app/url-shortener
COPY --from=builder /app/config /app/config

EXPOSE 8080

ENTRYPOINT ["/app/url-shortener"]
