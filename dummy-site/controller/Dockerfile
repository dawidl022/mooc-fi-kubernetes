FROM golang:1.19-alpine3.16 as build-stage
WORKDIR /usr/src/app
COPY . .
RUN go build -o controller


FROM alpine:3.16.2
COPY --from=build-stage /usr/src/app .
CMD [ "./controller" ]
