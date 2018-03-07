OWNER=sjeandeaux
REPO=github-cmd-go
SRC_DIR=github.com/$(OWNER)/$(REPO)
BUILD_VERSION=$(shell cat VERSION.txt)
NEXT_VERSION?=$(shell incrementor -position minor -version $(BUILD_VERSION))
#Default application or lambda
APPL?=associator

######## commom

PKGGOFILES=$(shell go list ./... | grep -v /vendor/)
CMD_TO_BUILD := ${sort ${dir ${wildcard ./cmd/*/}}}

GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_DIRTY?=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_DESCRIBE?=$(shell git describe --tags --always)
BUILD_TIME?=$(shell date +"%Y-%m-%dT%H:%M:%S")

LDFLAGS=-ldflags "\
          -X $(SRC_DIR)/information.Version=$(BUILD_VERSION) \
          -X $(SRC_DIR)/information.BuildTime=$(BUILD_TIME) \
          -X $(SRC_DIR)/information.GitCommit=$(GIT_COMMIT) \
          -X $(SRC_DIR)/information.GitDirty=$(GIT_DIRTY) \
          -X $(SRC_DIR)/information.GitDescribe=$(GIT_DESCRIBE)"


PWD=$(shell pwd)

GOOS?=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH?=amd64

define build-and-associate
	GOOS=$(1) GOARCH=$(2) go build $(LDFLAGS) -o ./target/$(1)-$(2)-${APPL} ./cmd/${APPL}
	GOOS=$(1) GOARCH=$(2) associator $(3) -name $(1)-$(2)-${APPL} -label $(1)-$(2)-${APPL} -content-type application/binary -owner $(OWNER) -repo $(REPO) -tag $(BUILD_VERSION)  -file ./target/$(1)-$(2)-${APPL}
endef

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-]+.*?:.*?## .*$$' ${MAKEFILE_LIST} | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[0;49;95m%-30s\033[0m %s\n", $$1, $$2}'


## If you have go on your wonderful laptop
.PHONY: clean
clean:
	@rm -rf ./target || true
	@mkdir ./target || true

.PHONY: test
test: fmt vet ## go test
	go test -cpu=2 -p=2 -race -v --short $(LDFLAGS) $(PKGGOFILES)

.PHONY: test-it-test
test-it-test: fmt vet ## go test with integration
	go test $(PKGGOFILES) -cpu=2 -p=2 -race  -v $(LDFLAGS)

.PHONY: test-cover
test-cover: fmt vet ## go test with coverage
	go test  $(PKGGOFILES) -cover -race -v $(LDFLAGS) -covermode=count -coverprofile=coverage.out

.PHONY: test-coverage
test-coverage: clean fmt vet ## for jenkins
	gocov test $(PKGGOFILES) --short -cpu=2 -p=2 -v $(LDFLAGS) | gocov-xml > ./target/coverage-test.xml

.PHONY: test-it-test-coverage
test-it-test-coverage: clean fmt vet ## for jenkins
	gocov test $(PKGGOFILES) -cpu=2 -p=2 -v $(LDFLAGS) | gocov-xml > ./target/coverage-test-it-test.xml

.PHONY: dependencies
dependencies: ## download the dependencies
	rm -rf Gopkg.lock vendor/
	dep ensure

.PHONY: build
build: clean fmt vet
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./target/${APPL} ./cmd/${APPL}

.PHONY: run
run: build ## Run command line
	go run $(LDFLAGS) ./cmd/${APPL}/main.go

.PHONY: build-lambda
build-lambda: build ## Build the zip for lambda
	zip --junk-paths ./target/${APPL}.zip ./target/${APPL}

.PHONY: run-lambda
run-lambda: build ## Run lambda through docker
	docker run --rm --env-file run-lambda-environments.properties -v ${PWD}/target:/var/task lambci/lambda:go1.x ${APPL}

.PHONY: fmt
fmt: ## go fmt on packages
	go fmt $(PKGGOFILES)

.PHONY: vet
vet: ## go vet on packages
	go vet $(PKGGOFILES)

.PHONY: lint
lint: ## go vet on packages
	golint -set_exit_status=true $(PKGGOFILES)

.PHONY: install
install: ## run 'go install' for each cmd
	@$(foreach dir,$(CMD_TO_BUILD),go install $(LDFLAGS) $(dir);)

.PHONY: tools
tools: ## install tools to develop
	go get -u github.com/golang/dep/cmd/dep
	go get github.com/sjeandeaux/github-cmd-go/cmd/associator
	go get github.com/sjeandeaux/github-cmd-go/cmd/incrementor
	go get -u github.com/golang/lint/golint
	go get github.com/axw/gocov/...
	go get github.com/AlekSi/gocov-xml

release-start: ## start release
	git flow release start $(BUILD_VERSION)

release: release-start ## create release
	git flow release finish -n -m "$(BUILD_VERSION)" $(BUILD_VERSION) 
	git checkout develop
	echo $(NEXT_VERSION) > VERSION.txt
	git add VERSION.txt
	git commit -m "bump the version"
	git push -u origin develop
	git push -u origin master
	git checkout master
	$(call build-and-associate,linux,amd64,-create)
	$(call build-and-associate,darwin,amd64)

release-finish-rollback:
	git tag -d $(BUILD_VERSION)
	git checkout master
	git reset --hard @{u}
	git checkout develop
	git reset --hard @{u}
