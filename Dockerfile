# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY go.mod ./
RUN go mod download

# throw err if no certs dir
COPY ./certs ./
COPY ./.env ./
COPY . ./
RUN mkdir -p ./log

ENV PORT=80

RUN CGO_ENABLED=0 GOOS=linux go build -o ./endpoint

EXPOSE 80

CMD ["/app/endpoint"]

