FROM golang:1.13.8-alpine3.11

RUN mkdir -p /go/src/messages
RUN apk add git

WORKDIR /go/src/messages

RUN go get github.com/tools/godep

# install dependencies
RUN mkdir -p /Godeps
COPY /Godeps /go/src/messages/Godeps
RUN godep restore

COPY . /go/src/messages
EXPOSE 8000

CMD ["go", "run", "messages.go"]
