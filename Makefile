pwd = $(shell pwd)

.PHONY: mac-cli server

all: mac-cli server

docker-cli: schema
	docker build -t fermayo/charger-cli -f cli/Dockerfile .

mac-cli: schema
	docker run --rm -v $(pwd):/go/src/github.com/fermayo/charger -w /go/src/github.com/fermayo/charger/cli -e GOOS=darwin -e GOARCH=amd64 golang:1.6 go build -v -o ../build/charger

server: schema
	docker build -t fermayo/charger-server -f server/Dockerfile .

schema: charger/charger.pb.go

charger/charger.pb.go: charger/charger.proto
	docker run -it --rm -v $(pwd):/go/src appcelerator/protoc -I charger/ charger/charger.proto --go_out=plugins=grpc:charger
