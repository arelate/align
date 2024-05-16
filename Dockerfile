FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o align \
    -ldflags="-s -w -X 'github.com/boggydigital/align/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

# adding align
FROM alpine:latest
COPY --from=build /go/src/app/align /usr/bin/align
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# https://en.wikipedia.org/wiki/Strategy_guide
EXPOSE 1981

# backups
VOLUME /usr/share/align/backups
# metadata
VOLUME /usr/share/align/metadata

ENTRYPOINT ["/usr/bin/align"]
CMD ["serve","-port", "1981", "-stderr"]

