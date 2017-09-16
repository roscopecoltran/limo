TEST?=./...

default: test

deps:
	glide install --strip-vendor

darwin:
	gox -verbose -os="darwin" -arch="amd64" -output="dist/{{.Dir}}" ./cmd/...

linux:
	gox -verbose -os="linux" -arch="amd64" -output="dist/{{.Dir}}" ./cmd/...

bin:
	@sh -c "$(CURDIR)/scripts/build.sh"

dev:
	@TF_DEV=1 sh -c "$(CURDIR)/scripts/build.sh"

test:
	go test $(TEST) $(TESTARGS) -timeout=10s

testrace:
	go test -race $(TEST) $(TESTARGS)

updatedeps:
	go get -d -v -p 2 ./...

.PHONY: bin default test updatedeps
