FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p logs && chmod 777 logs

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata libc6-compat

COPY --from=builder /app/main .
COPY --from=builder /app/logs ./logs

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser:appgroup logs

USER appuser

EXPOSE 3000

ENTRYPOINT ["./main"]