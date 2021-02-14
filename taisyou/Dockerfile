# build init
FROM golang:1.15-buster

COPY ./init /go/src/init

WORKDIR /go/src/init

ENV GO111MODULE on

RUN go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc

RUN GOOS=linux GOARCH=amd64 go build -o main main.go

# build phaul
FROM golang:1.15-buster

COPY ./phaul /go/src/phaul

WORKDIR /go/src/phaul

RUN GOOS=linux GOARCH=amd64 go build -o c/phaul-client c/c.go
RUN GOOS=linux GOARCH=amd64 go build -o s/phaul-server s/s.go

# main app
FROM ruby:2.6-stretch

WORKDIR /root/

COPY --from=0 /go/src/init/main .
COPY --from=1 /go/src/phaul/ ./phaul/

RUN echo "deb http://ftp.jp.debian.org/debian sid main" >> /etc/apt/sources.list \
    && apt update \
    && apt install -y libc6 -f -o APT::Immediate-Configure=0 \
    && apt install -y --no-install-recommends criu \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*


COPY zipkin-ruby-example /root/zipkin-ruby-example
COPY start.sh /root

RUN bundle install --gemfile=/root/zipkin-ruby-example/Gemfile
EXPOSE 80 8081 9000 50051
