FROM golang:1.23-alpine

RUN apk update &&  apk add git
RUN go install github.com/air-verse/air@latest
WORKDIR /app

CMD ["air", "-c", ".air.toml", "-d"]
