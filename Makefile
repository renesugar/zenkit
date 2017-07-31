ZENKIT_VERSION       := 1.3
GINKGO               := $(shell command -v ginkgo 2> /dev/null)
PACKAGE              := github.com/zenoss/zenkit
LOCAL_USER_ID        := $(shell id -u)
BUILD_IMG            := zenoss/zenkit-build:$(ZENKIT_VERSION)

ifndef IN_DOCKER
DOCKER_CMD           := docker run --rm -t \
							-v $(GOPATH)/src:/go/src:rw \
							-e LOCAL_USER_ID=$(LOCAL_USER_ID) \
							-e IN_DOCKER=1 \
							-w /go/src/$(PACKAGE) \
							$(BUILD_IMG)
else
DOCKER_CMD           :=
endif

.PHONY: default
default: test

.PHONY: test-containerized
test-containerized:
	@$(DOCKER_CMD) make test

.PHONY: test
ifndef GINKGO
test: test-containerized
else
test: COVERAGE_PROFILE := zenkit.coverprofile
test: COVERAGE_XML     := zenkit.coverage.xml
test:
	@ginkgo -r \
		-cover \
		-covermode=count \
		-tags="unit integration" \
		--skipPackage vendor
	@gocov convert $(COVERAGE_PROFILE) | gocov-xml > $(COVERAGE_XML)
endif

.PHONY: clean
clean:
	rm -rf zenkit.coverprofile zenkit.coverage.xml

local-dev:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/wadey/gocovmerge
	go get -u github.com/axw/gocov/gocov
	go get -u github.com/AlekSi/gocov-xml
