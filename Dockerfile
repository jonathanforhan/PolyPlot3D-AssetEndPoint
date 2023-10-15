# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./
COPY assets/. ./
RUN mkdir -p ./log

RUN CGO_ENABLED=0 GOOS=linux go build -o ./endpoint

EXPOSE 80

CMD ["/app/endpoint"]

