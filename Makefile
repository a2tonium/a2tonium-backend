#LOCAL_BIN:=$(CURDIR)/bin
#
#ifeq ($(wildcard $(LOCAL_BIN)),)
#$(shell mkdir "$(LOCAL_BIN)")
#endif
#
#.PHONY: all
#all: test build
#
#.PHONY: test
#test:
#	$(info Running tests...)
#	go test ./... -timeout 30s
#
#.PHONY: build
#build:
#	$(info Building...)
#	go build -o="$(LOCAL_BIN)" -ldflags="-s -w" ./cmd/ehed/
#
#.PHONY: build_debug
#build_debug:
#	$(info Building...)
#	go build -o="$(LOCAL_BIN)" ./cmd/
#
.PHONY: run
run:
	$(info Running...)
	go run ./cmd/a2tonium -config=local
	#go run ./cmd/loadtest-service -config=local
