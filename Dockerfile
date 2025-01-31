FROM golang:1.23.5 as builder
WORKDIR /shorten_url
ADD . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/shorten_url ./cmd

FROM alpine:3.16
WORKDIR /shorten_url
RUN apk update && \
    apk upgrade && \
    apk add --no-cache curl tzdata && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /shorten_url/cmd/shorten_url /shorten_url/shorten_url

EXPOSE 8080

CMD ["/shorten_url/shorten_url", "run"]
