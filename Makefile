# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = sympton-checker
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION?=?

# Symlink into GOPATH
GITHUB_USERNAME=bjornaer
BUILD_DIR=~/personal/${BINARY}
MAIN_DIR=./cmd/SymptomChecker
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: link clean test vet linux darwin windows

link:
	BUILD_DIR=${BUILD_DIR}; \
	BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
	CURRENT_DIR=${CURRENT_DIR}; \
	if [ "$${BUILD_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
	    echo "Fixing symlinks for build"; \
	    rm -f $${BUILD_DIR}; \
	    ln -s $${CURRENT_DIR} $${BUILD_DIR}; \
	fi

linux: 
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} ${MAIN_DIR} ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} ${MAIN_DIR} ; \
	cd - >/dev/null

windows:
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe ${MAIN_DIR} ; \
	cd - >/dev/null

run:
	cd ${BUILD_DIR}; \
	go run ${MAIN_DIR}

test:
	cd ${BUILD_DIR}; \
	go test -v -race ./...

vet:
	-cd ${BUILD_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/);

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

local:
	cd ${BUILD_DIR}; \
	go run ${MAIN_DIR} & npm run --prefix frontend start &

docker:
	cd ${BUILD_DIR}; \
	docker-compose up

.PHONY: link linux darwin windows test vet fmt clean run local docker