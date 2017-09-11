.PHONY: all deps build

# determine platform
ifeq (Darwin, $(findstring Darwin, $(shell uname -a)))
  MACHINE_LOCAL_PLATFORM := darwin
else
  MACHINE_LOCAL_PLATFORM := linux
endif	

# This Makefile is a simple example that demonstrates usual steps to build a binary that can be run in the same
# architecture that was compiled in. The "ldflags" in the build assure that any needed dependency is included in the
# binary and no external dependencies are needed to run the service.

KRAKEND_VERSION=$(shell git describe --always --long --dirty --tags)
BIN_NAME=krakend_searx_example_${KRAKEND_VERSION}

all: deps build

deps:
	go get "github.com/devopsfaith/krakend/config/viper"
	go get "github.com/devopsfaith/krakend/proxy"
	go get "github.com/devopsfaith/krakend/router/gin"
	go get "github.com/devopsfaith/krakend/logging/gologging"
	go get "github.com/devopsfaith/krakend/sd/etcd"

build:
	go build -a -ldflags="-X github.com/devopsfaith/krakend/core.KrakendVersion=${KRAKEND_VERSION}" -o ${BIN_NAME}
	@echo "You can now use ./${BIN_NAME}"

### GOLANG #####################################################################

fmt:
	go fmt ./...

install-deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

### GOLANG_GLIDE ###############################################################

GLIDE_EXEC						:= $(which glide)
GLIDE_CREATE_AUTO 				:= True
GLIDE_CREATE_AUTO_ARG 			:= "--non-interactive"

GLIDE_INSTALL_FORCE 			:= True
GLIDE_INSTALL_FORCE_ARG 		:= "--force"
GLIDE_INSTALL_STRIP_VENDOR 		:= True
GLIDE_INSTALL_STRIP_VENDOR_ARG 	:= "--strip-vendor"

glide-get:
	@if [ ! -f $(GLIDE_EXEC) ]; then go get -v github.com/Masterminds/glide; fi

glide-create:
	@if [ ! -f glide.yaml ]; then glide create $(GLIDE_CREATE_AUTO_ARG); fi

glide-install: glide-create
	glide install $(GLIDE_INSTALL_FORCE_ARG) $(GLIDE_INSTALL_STRIP_VENDOR_ARG)

glide-prepare: glide-get glide-create glide-install

### GOLANG_GOM #################################################################

gom-get:
	go get -v github.com/mattn/gom

gom-install:
	gom gen gomfile

### GOLANG_GOX #################################################################

gox-cross: glide-get  glide-get gox-darwin gox-linux gox-windows
# gox -verbose -os="linux darwin windows" -arch="amd64" -output="/shared/dist/{{.Dir}}/{{.Dir}}_{{.OS}}_{{.ARCH}}" $(glide novendor)

gox-all: glide-get  gox-local gox-cross

gox-local: glide-create
	gox -verbose -os="$(MACHINE_LOCAL_PLATFORM)" -arch="amd64" -output="./bin/{{.Dir}}" $(glide novendor)

gox-darwin: glide-get glide-create
	gox -verbose -os="darwin" -arch="amd64" -output="./dist/{{.Dir}}/{{.Dir}}_{{.OS}}_{{.ARCH}}" $(glide novendor)

gox-linux: glide-get glide-create
	gox -verbose -os="linux" -arch="amd64" -output="./dist/{{.Dir}}/{{.Dir}}_{{.OS}}_{{.ARCH}}" $(glide novendor)

gox-windows: glide-get glide-create
	gox -verbose -os="windows" -arch="amd64" -output="./dist/{{.Dir}}/{{.Dir}}_{{.OS}}_{{.ARCH}}" $(glide novendor)

### GOLANG_FIX #################################################################

logrus-fix:
	rm -fr vendor/github.com/Sirupsen
	find vendor -type f -exec sed -i 's/Sirupsen/sirupsen/g' {} +

### GOLANG_GENERATE ############################################################

generate: clean generate-models

generate-proto:
	protoc --gogofaster_out=. -Iproto -I$(GOPATH)/src proto/caffe2.proto

generate-models:
	go-bindata -nomemcopy -prefix builtin_models/ -pkg caffe2 -o builtin_models_static.go -ignore=.DS_Store  -ignore=README.md builtin_models/...

### GOLANG_CLEAN ###############################################################

clean-models:
	rm -fr builtin_models_static.go

clean-proto:
	rm -fr *pb.go

clean: clean-models

travis: install-deps glide-install logrus-fix generate
	echo "building..."
	go build