<<<<<<< HEAD
# --- Stage 1: Dependencies ---
FROM node:18-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

# --- Stage 2: Build ---
FROM node:18-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# --- Stage 3: Final ---
FROM node:18-alpine
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app .
CMD ["node", "index.js"]
=======
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
>>>>>>> 47244df147086597217fcdc9446abba65e38d84d
