GO_SRC := $(shell find . -type f -name "*.go")
LD_FLAGS := "-w -s"

start:
	docker-compose up --build -d
	docker-compose scale worker=3

restart:
	docker-compose up --no-deps --build -d worker

stop:
	docker-compose down

kbmb: $(GO_SRC)
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LD_FLAGS) -o $@ -a -installsuffix cgo .