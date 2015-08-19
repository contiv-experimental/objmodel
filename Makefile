
.PHONY: all build checks
TO_BUILD := ./tools/modelgen/ ./objdb/ ./objdb/client/ ./objdb/plugins/ ./objdb/plugins/etcdClient/ ./contivModel/ ./contivModel/cmExample/

all: generate test build

godep:
	@if [ -z "`which godep`" ]; then go get -v github.com/kr/godep; fi

vet:
	@(go tool | grep vet) || go get -v golang.org/x/tools/cmd/vet

checks: vet
	./checks "$(TO_BUILD)"

generate:
	cd tools/modelgen/generators && sh build.sh >templates.go && gofmt -w -s templates.go

build: godep checks
	godep go install -v ./...
	make clean

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

host-test: etcd build
	godep go test -v ./...
	make clean

reflex:
	# go get github.com/cespare/reflex
	reflex -r '.*\.go' make test
