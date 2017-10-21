FROM golang:1.9.1

RUN go get -u  github.com/todaychiji/ha/cli/gateway && \
    go install github.com/todaychiji/ha/cli/gateway

ENTRYPOINT gateway