name=ansible
component=manager
version=2.5.0
release=3
docker_registry=192.168.1.100:5000
current_dir = $(shell pwd)
project_dir = github.com/gzsunrun/ansible-manager

.PHONY: bin rpm docker
all: bin

bindata:
	@echo ">> bindata .."
	@go-bindata -o=asset/asset.go -pkg=asset public/...
# build binary
bin:
	rm -rf ../../vendor/go.etcd.io/etcd/vendor/golang.org/x/net/trace
	go build --tags consul -ldflags "-s -w -X main._VERSION_=$(version)-$(release)" -o ./bin/$(name)-$(component)
# build rpm package
rpm: bin
	rm -rf ~/rpmbuild/BUILD/$(name)-$(component)/
	rm -f ~/rpmbuild/SPECS/$(name)-$(component).spec
	mkdir -p ~/rpmbuild/BUILD/$(name)-$(component)/etc/systemd/system/
	mkdir -p ~/rpmbuild/BUILD/$(name)-$(component)/usr/local/bin/
	cp ./build/$(name)-$(component).service ~/rpmbuild/BUILD/$(name)-$(component)/etc/systemd/system/
	cp ../../bin/$(name)-$(component) ~/rpmbuild/BUILD/$(name)-$(component)/usr/local/bin/
	mkdir -p ~/rpmbuild/SPECS
	cp ./build/$(name)-$(component).spec ~/rpmbuild/SPECS/
	sed -i 's/^Version:.*/Version:$(version)/' ~/rpmbuild/SPECS/$(name)-$(component).spec
	sed -i 's/^Release:.*/Release:$(release).el7/' ~/rpmbuild/SPECS/$(name)-$(component).spec
	rpmbuild -bb ~/rpmbuild/SPECS/$(name)-$(component).spec --define "_rpmdir $(current_dir)/../../rpm/"
# build docker image
docker:
	sed -i 's/^FROM.*/FROM $(docker_registry)\/library\/alpine:3.8-ansible/' ./Dockerfile 
	docker run --rm -e GOROOT=/usr/local/go -e GOPATH=/gopath -v $(current_dir)/:/gopath/src/$(project_dir) -w /gopath/src/$(project_dir)  $(docker_registry)/library/golang:1.10-alpine-ssh  make
	cp -f ./bin/$(name)-$(component) .
	docker build -t $(name)-$(component):$(version)-$(release) .
	rm -f ./bin/$(name)-$(component)
# push docker image
docker_push:
	docker tag $(name)-$(component):$(version)-$(release) $(docker_registry)/gzsunrun/$(name)-$(component):$(version)-$(release)
	docker tag $(name)-$(component):$(version)-$(release) $(docker_registry)/gzsunrun/$(name)-$(component):latest
	docker push $(docker_registry)/gzsunrun/$(name)-$(component):$(version)-$(release)
	docker push $(docker_registry)/gzsunrun/$(name)-$(component):latest
