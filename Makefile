
.PHONY: all build checks
TO_BUILD := ./tools/modelgen/ ./objdb/ ./objdb/client/ ./objdb/plugins/ ./objdb/plugins/etcdClient/ ./contivModel/ ./contivModel/cmExample/

all: build run-test

godep:
	@if [ -z "`which godep`" ]; then go get -v github.com/kr/godep; fi

vet:
	@(go tool | grep vet) || go get -v golang.org/x/tools/cmd/vet

checks: vet
	./checks "$(TO_BUILD)"

generator:
	cd tools/modelgen/generators && bash build.sh >templates.go && gofmt -w -s templates.go

generate:
	cd contivModel && bash generate.sh

build: godep checks generator
	make clean
	godep go install -v ./...
	make generate

etcd:
	pkill etcd || exit 0
	rm -rf /tmp/etcd
	etcd --force-new-cluster --data-dir /tmp/etcd &

build-docker:
	docker build -t objmodel .

test: build-docker
	docker run --rm -v ${PWD}:/gopath/src/github.com/contiv/objmodel objmodel

clean:
	rm -rf Godeps/_workspace/pkg

run-test: 
	godep go test -v ./...

host-test: etcd build
	make run-test
	make clean

reflex:
	# go get github.com/cespare/reflex
	reflex -r '.*\.go' -R tools/modelgen/generators/templates.go make test
