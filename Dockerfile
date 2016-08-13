FROM alpine:latest
MAINTAINER  Piotr Kowalczuk <p.kowalczuk.priv@gmail.com>
COPY ./bin /usr/local/bin/
COPY ./docker-entrypoint.sh /

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["mnemosyned"]