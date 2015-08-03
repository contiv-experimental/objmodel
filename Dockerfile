FROM debian:latest

ENV DEBIAN_FRONTEND noninteractive
ENV ETCD_VER v2.0.10
ENV ETCD_FILE etcd-${ETCD_VER}-linux-amd64

RUN apt-get update && apt-get install build-essential curl git mercurial -y

RUN curl -sL https://github.com/coreos/etcd/releases/download/${ETCD_VER}/${ETCD_FILE}.tar.gz | tar vxz -C /tmp && cp /tmp/${ETCD_FILE}/etcd* /usr/bin

RUN curl https://storage.googleapis.com/golang/go1.4.2.src.tar.gz | tar -C /usr/local -xvz
RUN cd /usr/local/go/src; for i in linux; do GOOS=$i GOARCH=amd64 ./make.bash; done
ENV GOBIN /gobin
ENV GOPATH /gopath
RUN mkdir ${GOPATH} ${GOBIN}
ENV PATH /usr/local/go/bin:${GOBIN}:${PATH}

ENV FULLPATH /gopath/src/github.com/contiv/objmodel

WORKDIR ${FULLPATH}

CMD make host-test
