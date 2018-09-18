FROM grpc/go:latest as grpc
COPY *.proto /prime/
WORKDIR /prime
RUN protoc -I /prime /prime/example.proto --go_out=plugins=grpc:.
RUN cat /prime/example.pb.go

FROM golang:1.11.0 as builder
ADD ./ /go/src/github.com/mattpaletta/prime-blockchain
WORKDIR /go/src/github.com/mattpaletta/prime-blockchain
COPY --from=grpc /prime/example.pb.go /go/src/github.com/mattpaletta/prime-blockchain/blockchain/example.pb.go
RUN go get -v ./...
RUN go install github.com/mattpaletta/prime-blockchain
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags netgo -installsuffix netgo -o app .

FROM alpine:latest
RUN apk add --no-cache bash iputils vim
ENV PATH=/bin

WORKDIR /root/
COPY --from=builder /go/src/github.com/mattpaletta/prime-blockchain/app .
ENTRYPOINT ["./app"]
