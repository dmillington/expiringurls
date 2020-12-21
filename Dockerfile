FROM golang:1.15.6 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev

RUN go build -o /server

EXPOSE 8000

CMD ["/server"]
