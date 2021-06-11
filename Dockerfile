FROM golang:1.16 AS build
WORKDIR /go/src/github.com/mkorenkov/sha256msg
COPY . .
RUN cd /go/src/github.com/mkorenkov/sha256msg && CGO_ENABLED=0 GO111MODULE=off GOOS=linux go build -o bin/sha256msg cmd/sha256msg/main.go

FROM alpine:latest
# explicitly set user/group IDs
RUN addgroup -S -g 998 sha256msg && \
    adduser -S -h /srv/sha256msg -u 998 -G sha256msg sha256msg && \
    apk add --update \
        bash \
        bind-tools \
        ca-certificates \
        su-exec \
        tzdata
COPY --from=build /go/src/github.com/mkorenkov/sha256msg/bin/sha256msg /bin/sha256msg
ENTRYPOINT ["su-exec", "sha256msg"]
CMD ["/bin/sha256msg"]
