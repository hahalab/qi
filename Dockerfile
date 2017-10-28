FROM golang:1.9.1

ADD . /go/src/github.com/hahalab/qi
RUN go install github.com/hahalab/qi/cli/qi

ENTRYPOINT ["qi", "gateway"]