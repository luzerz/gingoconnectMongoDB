FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . ./


RUN go build -o /main .

EXPOSE 80

ENTRYPOINT [ "/main" ]


# FROM alpine:3.13
# RUN apk update && apk --no-cache add curl && apk --no-cache add tzdata
# ADD main /main

