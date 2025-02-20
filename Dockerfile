FROM golang:1.23.5-alpine as builder
WORKDIR /shorten_url
ADD . .

RUN apk add --no-cache gcc musl-dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/shorten_url ./cmd

FROM alpine:3.16
WORKDIR /shorten_url
RUN apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

COPY --from=builder /shorten_url/cmd/shorten_url /shorten_url/shorten_url

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/v1 || exit 1

CMD ["/shorten_url/shorten_url", "run"]
