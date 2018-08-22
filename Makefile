unexport GOBIN

NAME         	?= ansible-manager
GO           	?= go
SPECE_DIR       ?= ./_rpmbuild
RPMBUILD_DIR	?= /root/rpmbuild
VERSION = 2.4.1
RELEASE = 4
DOCKERREPO = 192.168.1.100:5000

all: pkg build


uppkg:
	@echo ">> update pkg..."
	glide --home /root/.glide/ansible-manager/  up

pkg:
	@echo ">> get pkg..."
	glide --home /root/.glide/ansible-manager/  install

build:
	@echo ">> building code..."
	go build -ldflags "-s -w" -o $(NAME)

release:
	git tag release-$(VERSION)-$(RELEASE)

dockerbuild:
	@echo ">> build docker image"
	docker build -t $(DOCKERREPO)/gzsunrun/$(NAME):master .

pimage:
	@echo ">> push release docker image"
	docker tag $(DOCKERREPO)/gzsunrun/$(NAME):master $(DOCKERREPO)/gzsunrun/$(NAME):$(VERSION)-$(RELEASE)
	docker tag $(DOCKERREPO)/gzsunrun/$(NAME):master $(DOCKERREPO)/gzsunrun/$(NAME):latest
	docker push $(DOCKERREPO)/gzsunrun/$(NAME):$(VERSION)-$(RELEASE)
	docker push $(DOCKERREPO)/gzsunrun/$(NAME):latest

pdevimage:
	@echo ">> push dev docker image"
	docker push  $(DOCKERREPO)/gzsunrun/$(NAME):master