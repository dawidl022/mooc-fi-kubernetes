FROM golang:1.18-alpine3.15 AS build-stage

WORKDIR /usr/src/app

COPY . .

RUN go build -o todo-server


FROM alpine:3.15.4

RUN adduser -D appuser

COPY --from=build-stage /usr/src/app/todo-server ./todo-server

USER appuser

CMD ["./todo-server"]
