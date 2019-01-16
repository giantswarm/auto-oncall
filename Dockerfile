FROM alpine:3.8
RUN apk add --no-cache ca-certificates

ADD ./oncall /oncall

EXPOSE 8000
ENTRYPOINT ["/oncall"]
