FROM golang:1.8

WORKDIR /go/src/github.com/incu6us/asterisk-ami-api
COPY . .
RUN curl https://glide.sh/get | sh
RUN glide i
RUN go build -o asterisk-ami-api main.go

EXPOSE 3000

CMD ["asterisk-ami-api", "-conf", "api.conf"]
