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