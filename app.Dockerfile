FROM golang:1.21 as builder

WORKDIR /usr/src/app

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download && go mod verify

COPY ./src .
RUN go build -o /usr/local/bin/app ./main.go

###
FROM ubuntu:jammy

WORKDIR /usr/local/app

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates

COPY --from=builder /usr/local/bin/app ./
COPY --from=builder /usr/src/app/templates ./templates

CMD ["./app"]
