FROM golang:1.7

RUN mkdir -p /go/src/github.com/fermayo/charger
COPY . /go/src/github.com/fermayo/charger

WORKDIR /go/src/github.com/fermayo/charger/cli
RUN go get -d -v && go install -v

ENTRYPOINT ["cli"]
