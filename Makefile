
TO_BUILD := ./tools/modelgen/ ./objdb/ ./objdb/objdbClient/ ./objdb/plugins/ ./objdb/plugins/etcdClient/ ./contivModel/ ./contivModel/cmExample/

all: test binaries

get:
	go get -v ./...

build: get
	go install -v ./...

etcd:
	pkill etcd || exit 0
	rm -rf /tmp/etcd
	etcd --force-new-cluster --data-dir /tmp/etcd &

build-docker:
	docker build -t objmodel .

test: build-docker
	docker run --rm -v ${PWD}:/gopath/src/github.com/contiv/objmodel objmodel

host-test: etcd build
	PATH=${PWD}/bin:${PATH} go test -v ./...
	rm -rf pkg

reflex:
	reflex -r '.*\.go' make test
