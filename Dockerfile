
FROM golang
MAINTAINER  Piotr Kowalczuk <p.kowalczuk.priv@gmail.com>

ADD . /go/src/github.com/piotrkowalczuk/mnemosyne
WORKDIR /go/src/github.com/piotrkowalczuk/mnemosyne

RUN go get ./...
RUN go install ./cmd/mnemosyned
ENTRYPOINT /go/bin/mnemosyned
EXPOSE 8080