SERVICE="lvodQuery"

COMMIT_SHA1	:= $(shell git rev-parse HEAD)
TAG_VERSION	:= $(shell git describe --abbrev=0 --tags)
BUILD_TIME	:= $(shell date "+%F %T")
LDFLAGS		:= "-w -s \
		-X 'main.COMMIT_SHA1=1$(COMMIT_SHA1)' \
		-X 'main.VERSION=$(TAG_VERSION)' \
		-X 'main.BUILD_TIME=$(BUILD_TIME)' \
		"
DB_ACCOUNT 	:= "root"
DB_PASSWORD := "123456"

ifdef TAG_VERSION
	TAG_VERSION := 'v$(TAG_VERSION)'
endif

ifndef TAG_VERSION
	TAG_VERSION := 'v$(shell cat VERSION)'
endif

SHORT_SHA1 := $(shell git rev-parse --short HEAD)

PACKAGEFILE := "etc" "lvodQuery" "Makefile" "lvodQueryConsole.sh" "install.sh"

dbinit:
	mysql -u$(DB_ACCOUNT) -p$(DB_PASSWORD) < etc/database/authority.sql

build:
	go build -ldflags $(LDFLAGS) -o $(SERVICE)

packages = $(shell go list ./...)
test:
	go test -v ${packages}

lint:
	@golangci-lint --version
	@golangci-lint run --deadline=5m ./...

linux: 
	GOOS=linux go build -ldflags $(LDFLAGS) -o $(SERVICE)

images: linux
	docker build -f Dockerfile -t $(SERVICE):$(TAG_VERSION)-$(shell git rev-parse --short HEAD) .

ALL: clean images

pack: linux
	tar -zvcf $(SERVICE)-$(TAG_VERSION)-$(SHORT_SHA1).tgz $(PACKAGEFILE)

.PHONY: clean ALL build build_proto dbinit install