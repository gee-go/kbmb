FROM iron/go:dev

ENV GOPATH /src/go
WORKDIR /src/go/src/github.com/gee-go/kbmb
ADD . /src/go/src/github.com/gee-go/kbmb
