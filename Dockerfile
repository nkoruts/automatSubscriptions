# build stage
FROM golang:1.25.5 AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# runtime stage
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/app /app/app

EXPOSE 9091
CMD ["/app/app"]