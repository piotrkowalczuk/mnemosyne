
FROM golang
MAINTAINER  Piotr Kowalczuk <p.kowalczuk.priv@gmail.com>

ADD . /go/src/github.com/piotrkowalczuk/mnemosyne

WORKDIR /go/src/github.com/piotrkowalczuk/mnemosyne

RUN make get
RUN go install github.com/piotrkowalczuk/mnemosyne/cmd/mnemosyned
RUN rm -rf /go/src

EXPOSE 8080

ENTRYPOINT ["/go/bin/mnemosyned"]
CMD ["-host=0.0.0.0", "-namespace=mnemosyne"]