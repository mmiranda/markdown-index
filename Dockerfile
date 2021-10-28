FROM golang:1.17.2 AS builder
WORKDIR /go/src/github.com/mmiranda/markdown-index
RUN go get -d -v golang.org/x/net/html
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgop .

FROM alpine:latest
LABEL org.opencontainers.image.source https://github.com/mmiranda/markdown-index
LABEL org.opencontainers.image.description "Tool to generate a global Markdown Summary Index based on other markdown files"
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/mmiranda/markdown-index/markdown-index .

VOLUME /data

ENTRYPOINT ["/app/markdown-index", "--directory", "/data"]
CMD ["--directory", "/data"]