# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY modules ./modules
COPY cmd ./cmd

# Build

RUN ls -a
RUN CGO_ENABLED=0 GOOS=linux go build -o /.bin/basic ./cmd/main.go

EXPOSE 8080

# Run
CMD ["/.bin/basic"]
