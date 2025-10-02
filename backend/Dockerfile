FROM golang:1.21-alpine AS builder

WORKDIR /app

# Dependencies installieren
COPY go.mod go.sum* ./
RUN go mod download

# Code kopieren und bauen
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Production Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
