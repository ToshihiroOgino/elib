FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && \
  apk update && \
  apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1

COPY . .
RUN go build -o main .

FROM alpine:latest AS runtime

WORKDIR /app

COPY --from=builder /app/main /app/go.mod /app/go.sum ./
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

CMD ["/app/main"]
