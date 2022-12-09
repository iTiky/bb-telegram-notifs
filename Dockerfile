FROM golang:1.19 AS builder
WORKDIR /go/src/github.com/itiky/bb-telegram-notifs/
COPY . /go/src/github.com/itiky/bb-telegram-notifs/
RUN make build

FROM ubuntu:22.04
RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/itiky/bb-telegram-notifs/bbtt ./
CMD [ "./bbtt" ]
