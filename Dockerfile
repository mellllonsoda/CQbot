# --- Stage 1: Build ---
FROM golang:1.25-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app

# --- Stage 2: Final ---
FROM alpine:latest

RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

WORKDIR /root/

COPY --from=builder /go-app .

COPY quotes.json .
COPY keywords.json .

CMD ["./go-app"]