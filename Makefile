ROOT := agora.io/rtm-sdk
NAME := rtm-sdk
VERSION ?= $(shell git describe --tags --always --dirty)
DOCKER_LABELS ?= git-describe="$(shell date -u +v%Y%m%d)-$(shell git describe --tags --always --dirty)"

export SHELLOPTS := errexit
export GOFLAGS := -mod=vendor

REGISTRY ?= hub.agoralab.co/adc/
BASE_REGISTRY ?=

APP_ID ?=
USER_ID ?=

GO_VERSION ?= 1.16.4
BIN_DIR := $(GOPATH)/bin
CMD_DIR := ./cmd
OUTPUT_DIR := ./bin
BUILD_DIR := ./build
.PHONY: rtm-swig wrapper build build-linux
rtm-swig:
	@swig -go -cgo -c++ -soname libagora_rtm_sdk.so -intgosize 64 ./internal/rtmlib/lib.i
wrapper:
	@docker run --rm -t																	\
		-v $(PWD):/go/src/$(ROOT)                                     \
   		-w /go/src/$(ROOT)/internal/rtmlib                                                             \
   	  	gcc:10                       \
        	/bin/bash -c ' g++ -std=c++11  -shared  -fPIC -I.  \
        	-o libwraper.so  lib_wrap.cxx  -L. -lagora_rtm_sdk '  \



build: build-linux

build-linux:
	@docker run --rm -t                                                                \
	  -v $(PWD):/go/src/$(ROOT)                                                        \
	  -w /go/src/$(ROOT)                                                               \
	  -e GOOS=linux                                                                    \
	  -e GOARCH=amd64                                                                  \
	  -e GOPATH=/go                                                                    \
	  -e GOFLAGS=$(GOFLAGS)   	                                                       \
	  -e SHELLOPTS=$(SHELLOPTS)                                                        \
	   $(BASE_REGISTRY)golang:$(GO_VERSION)                       \
	    /bin/bash -c '                                    								\
	      	go build  -v -o $(OUTPUT_DIR)/$(NAME)					\
          		 -ldflags "-s -w"													\
          	 $(CMD_DIR)/'                                                    			\

container: build-linux
	@echo ">> building image"
	@docker build -t $(REGISTRY)$(NAME):$(VERSION) --label $(DOCKER_LABELS)  -f $(BUILD_DIR)/Dockerfile .

push: container
	@echo ">> pushing admin image"
	@docker push $(REGISTRY)$(NAME):$(VERSION)
run: container
	@echo ">> running image"
	@ docker run --rm -ti $(REGISTRY)$(NAME):$(VERSION)  -APP_ID=$(APP_ID) -USER_ID=$(USER_ID)