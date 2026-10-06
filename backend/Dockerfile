FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/notes .

FROM alpine:3.20

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

COPY --from=builder /out/notes /app/notes

RUN chown -R app:app /app

USER app

EXPOSE 8000

HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8000/health || exit 1

ENTRYPOINT ["/app/notes"]
