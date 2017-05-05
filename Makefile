build:
	go build

deps:
	go get -u github.com/aws/aws-sdk-go/...
	go get github.com/op/go-logging

