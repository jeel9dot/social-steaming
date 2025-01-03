FROM golang:1.23.4-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o social-stream .

FROM alpine:latest
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=builder /app/social-stream .
COPY .env.docker /app/.env
EXPOSE 50051
ENTRYPOINT ["./social-stream"]
CMD ["grpc"]
