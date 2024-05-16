FROM golang:alpine as build
RUN apk add --no-cache --update git
ADD . /go/src/app
WORKDIR /go/src/app
RUN go get ./...
RUN go build \
    -a -tags timetzdata \
    -o boilerplate \
    -ldflags="-s -w -X 'github.com/boggydigital/boilerplate/cli.GitTag=`git describe --tags --abbrev=0`'" \
    main.go

# adding boilerplate
FROM alpine:latest
COPY --from=build /go/src/app/boilerplate /usr/bin/boilerplate
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 1234

# backups
VOLUME /usr/share/boilerplate/backups
# metadata
VOLUME /usr/share/boilerplate/metadata

ENTRYPOINT ["/usr/bin/boilerplate"]
CMD ["serve","-port", "1234", "-stderr"]

