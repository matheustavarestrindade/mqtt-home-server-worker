# ---- Build Stage ----
FROM golang:1.24.4-alpine AS builder

# Set the working directory.
WORKDIR /app

# Copy and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the Go application into a static binary.
RUN CGO_ENABLED=0 go build -o /server cmd/server/main.go

# ---- Final Stage ----
FROM scratch

ARG BROKER_URL
ARG CLIENT_ID
ARG DATABASE_URL
ARG TELEGRAM_BOT_TOKEN
ARG TELEGRAM_CHAT_IDS
ARG MQTT_CLIENT_ID


ENV MQTT_CLIENT_ID=$MQTT_CLIENT_ID
ENV BROKER_URL=$BROKER_URL
ENV CLIENT_ID=$CLIENT_ID
ENV DATABASE_URL=$DATABASE_URL
ENV TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN
ENV TELEGRAM_CHAT_IDS=$TELEGRAM_CHAT_IDS


COPY certs /certs
COPY --from=builder /server /server

CMD ["/server"]
