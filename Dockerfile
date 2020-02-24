FROM golang:1.13.8-alpine3.11

WORKDIR /go/src/messages
COPY . .
EXPOSE 8000

RUN apk add git
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["go", "run", "messages.go"]
