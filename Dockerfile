FROM golang:1.17.2 AS builder
WORKDIR /go/src/github.com/mmiranda/markdown-index
RUN go get -d -v golang.org/x/net/html
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgop .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/mmiranda/markdown-index/markdown-index .
RUN pwd && ls

VOLUME /data

CMD ["/app/markdown-index"]
