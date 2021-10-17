FROM golang:1.17

RUN mkdir -p /app
RUN apt update
RUN apt install -y ffmpeg bc

WORKDIR /app

COPY . /app

# RUN go mod download

# RUN go build -o /app_exe

CMD bash -c create-vod-hls.sh ./tmp/sample.mp4