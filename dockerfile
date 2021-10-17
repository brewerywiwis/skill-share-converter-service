FROM golang:1.17

RUN mkdir -p /app
RUN mkdir -p /app/tmp
RUN apt update
RUN apt install -y ffmpeg bc

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o /app_exe

CMD /app_exe