DIRS=$(shell go list -f {{.Dir}} ./...)

# Only list test and build dependencies
# Standard dependencies are installed via go get
DEPEND=\
	golang.org/x/lint/golint \
	github.com/golangci/golangci-lint/cmd/golangci-lint \
	golang.org/x/tools/cmd/cover \
	golang.org/x/tools/cmd/goimports \
	gopkg.in/src-d/go-kallax.v1/... \
	github.com/gogo/protobuf/protoc-gen-gogo \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
	github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
	github.com/go-bindata/go-bindata/...


templates:
	go run tools/gen-tmpls/*.go

protoc-plugin:
	go install ./tools/protoc-gen-protoc-plugin

build:
	go build -o bex main.go

cleanup:
	rm bex || true
	rm -r a || true

precommit:	cleanup	templates

.PHONY: depend genbindata build templates  cleanup precommit protoc-plugin