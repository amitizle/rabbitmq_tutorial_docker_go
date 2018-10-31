# Step 1 compile
FROM golang:alpine as builder
COPY . $GOPATH/src/amqpexample/client/
WORKDIR $GOPATH/src/amqpexample/client/
RUN apk --update add git ca-certificates && \
		go get -d -v && \
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/client

# Step 2 build a small image
FROM scratch
COPY --from=builder /go/bin/client /go/bin/client
CMD ["/go/bin/client"]
