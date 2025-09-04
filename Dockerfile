FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app .

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/app .

COPY inspector_linux_amd64 .
RUN chmod +x inspector_linux_amd64 app

EXPOSE 64333

CMD ["./app", "-rtp=0.95"]
