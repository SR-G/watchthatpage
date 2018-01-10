SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=watchthatpage
PWD := $(shell pwd)

VERSION=1.0.0-SNAPSHOT
BUILD_TIME=$(date "%FT%T%z")

LDFLAGS=-ldflags "-d -s -w -X tensin.org/watchthatpage/core.Build=`git rev-parse HEAD`" -a -tags netgo -installsuffix netgo
PACKAGE=tensin.org/watchthatpage
ifeq ($(shell hostname),jupiter)
	DOCKER_IMAGE="tensin-app-golang"
else
	DOCKER_IMAGE="library/golang"
endif

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o bin/${BINARY} ${PACKAGE}

.PHONY: install clean deploy run 

build:
	time go install ${PACKAGE}

install:
	GOARCH=amd64 GOOS=windows go install ${PACKAGE}
	GOARCH=amd64 GOOS=linux go install ${LDFLAGS} ${PACKAGE}

clean:
	-@rm -f bin/cache/* 2>/dev/null || true
	-@rm -f bin/${BINARY} || true
	-@rm -f ${BINARY}-${VERSION}.zip 2>/dev/null || true

deploy:
	cp bin/watchthatpage* /home/applications/watchthatpage/
	cp -Rp resources/templates /home/applications/watchthatpage/

distribution: install
	-@mkdir /go/bin/linux/ || true
	mv /go/bin/watchthatpage /go/bin/linux/
	cp /go/resources/config/watchthatpage.json /go/bin/linux/
	cp /go/resources/config/watchthatpage.json /go/bin/windows_amd64/
	cp -Rp /go/resources/templates/ /go/bin/linux/
	cp -Rp /go/resources/templates/ /go/bin/windows_amd64/
	cd /go/bin/ ; zip -r -9 ${BINARY}-${VERSION}.zip ./linux/ ; zip -r -9 ${BINARY}-${VERSION}.zip ./windows_amd64/

run:
	bin/watchthatpage grab

init:
	[ ! -f bin/glide ] && glide.sh/get | sh
	glide update
	glide install
	
test:
	go test -v tensin.org/watchthatpage/core

docker:
	docker run --rm -it -v ${PWD}:/go ${DOCKER_IMAGE} /bin/bash
