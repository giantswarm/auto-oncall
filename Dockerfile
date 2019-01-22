FROM alpine:3.8
RUN apk add --no-cache ca-certificates

ADD ./auto-oncall /auto-oncall

EXPOSE 8000
ENTRYPOINT ["/auto-oncall"]
