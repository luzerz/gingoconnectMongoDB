FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

EXPOSE 8080

ENTRYPOINT [ "/main" ]


# FROM alpine:3.13
# RUN apk update && apk --no-cache add curl && apk --no-cache add tzdata
# ADD main /main

