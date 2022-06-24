FROM golang:1.18-alpine3.15 as build-stage

WORKDIR /usr/src/app

COPY . .

RUN go build -o current-log


FROM alpine:3.15.4

RUN adduser -D appuser

WORKDIR /usr/src/app

COPY --from=build-stage /usr/src/app .

USER appuser

RUN ["./current-log"]
