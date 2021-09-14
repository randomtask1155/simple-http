FROM golang:1.16 AS builder
LABEL maintainer="Daniel Lynch <danplynch@gmail.com>"
RUN mkdir -p /go/src/github.com/randomtask1155/simple-http
ENV GOPATH=/go
WORKDIR $GOPATH/src/github.com/randomtask1155/simple-http
ENV PATH=$GOPATH/bin:$PATH
ENV PORT=8080
ADD . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o simple-http .
RUN pwd

FROM alpine:latest as certs
RUN apk --update add ca-certificates


FROM scratch
COPY --from=builder /go/src/github.com/randomtask1155/simple-http/simple-http /go/bin/simple-http
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8080
ENTRYPOINT ["/go/bin/simple-http"]

