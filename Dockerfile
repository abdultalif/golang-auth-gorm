FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env.production .env.production 
COPY .env.production .env            

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata libc6-compat

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 3000

ENTRYPOINT ["./main"]