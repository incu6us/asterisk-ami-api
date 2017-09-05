FROM golang:1.8-alpine

WORKDIR /go/src/github.com/incu6us/asterisk-ami-api
COPY . .
RUN go get github.com/Masterminds/glide
RUN glide i
RUN go build -o asterisk-ami-api main.go
RUN rm -rf vendor

EXPOSE 3000

CMD ["./asterisk-ami-api", "-conf", "api.conf"]
