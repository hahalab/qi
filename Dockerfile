FROM golang:1.9.1

ADD . /go/src/github.com/todaychiji/ha
RUN go install github.com/todaychiji/ha/cli/up

ENTRYPOINT gateway gateway