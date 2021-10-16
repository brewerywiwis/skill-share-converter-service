FROM golang:1.17

RUN mkdir -p /app

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o /app_exe

CMD /app_exe $QUEUE_NAME