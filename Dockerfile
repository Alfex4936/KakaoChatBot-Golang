# To use git command for golang-mysql-driver
# FROM ubuntu:latest
# RUN apt-get -y update
# RUN apt-get -y install git

### Builder
FROM golang:1.16.0-alpine as builder
# RUN apk update && apk add git && apk add ca-certificates

WORKDIR $GOPATH/src/chatbot

# golang-mysql-driver requires git command
RUN set -ex && apk add --no-cache --virtual git

ADD . .

# Install all dependencies
RUN go get -d -v ./...

# where modules save /go/pkg/mod/github.com/
COPY soup.go /go/pkg/mod/github.com/anaskhan96/soup@v1.2.4

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o /bin/chatbot main/main.go

RUN chmod +x /bin/chatbot

ENV GO_MYSQL "ID:PASSWD@tcp(RDS_SERVER)/notices"
ENV GO111MODULE "auto"
EXPOSE 8008

### Make executable image
FROM scratch

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/chatbot /bin/chatbot

CMD [ "/bin/chatbot" ]