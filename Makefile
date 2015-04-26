NAME=auditor
VERSION=$(shell cat VERSION)

dev:
	docker build -f Dockerfile.dev -t $(NAME):dev .
	docker run --rm --net host \
		-v /var/run/docker.sock:/tmp/docker.sock \
		$(NAME):dev /bin/auditor logstash:

build:
	mkdir -p build
	docker build -t $(NAME):$(VERSION) .
	docker save $(NAME):$(VERSION) | gzip -9 > build/$(NAME)_$(VERSION).tgz
