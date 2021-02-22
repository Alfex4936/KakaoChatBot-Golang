# To use git command for golang-mysql-driver
FROM ubuntu:latest
RUN apt-get -y update
RUN apt-get -y install git

### Builder
FROM golang:1.16.0-alpine as builder
# RUN apk update && apk add git && apk add ca-certificates

WORKDIR $GOPATH/src/chatbot

ADD . .

# Install all dependencies
RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w -s' -o /bin/chatbot main/main.go

RUN chmod +x /bin/chatbot

ENV GO_MYSQL=ID:PASSWD@tcp(RDS_SERVER)/notices
ENV GO111MODULE="auto"
EXPOSE 8008

### Make executable image
FROM scratch

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /bin/chatbot /bin/chatbot

CMD [ "/bin/chatbot" ]