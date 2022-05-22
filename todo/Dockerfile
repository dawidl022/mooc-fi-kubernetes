FROM golang:1.18-alpine3.15 AS build-stage

WORKDIR /usr/src/app

COPY . .

RUN go build -o todo-server


FROM scratch

COPY --from=build-stage /usr/src/app .

CMD ["./todo-server"]
