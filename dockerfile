FROM golang:1.17-alpine

RUN mkdir -p /app
RUN mkdir -p /app/tmp
RUN apk update
RUN apk add --no-cache ffmpeg
RUN apk add --no-cache bc
RUN apk add --no-cache bash
WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o /app_exe

RUN mkdir -p /app/tmp
CMD /app_exe