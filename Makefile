pwd = $(shell pwd)

.PHONY: cli server

all: cli server

cli: schema
	docker build -t fermayo/charger-cli -f cli/Dockerfile .

server: schema
	docker build -t fermayo/charger-server -f server/Dockerfile .

schema: charger/charger.pb.go

charger/charger.pb.go: charger/charger.proto
	docker run -it --rm -v $(pwd):/go/src appcelerator/protoc -I charger/ charger/charger.proto --go_out=plugins=grpc:charger
