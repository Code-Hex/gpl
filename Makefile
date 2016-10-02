SUB_BIN   = sub_bin
PWD       = $(shell pwd)
REPO      = $(shell basename $(PWD))
CMDPATH   = ./cmd/$(REPO)/
GOVERSION = $(shell go version)
GOOS      = $(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH    = $(word 2,$(subst /, ,$(lastword $(GOVERSION))))
HAS_GLIDE = $(shell which glide)
XC_OS     = "darwin linux windows"
XC_ARCH   = "386 amd64"
VERSION   = $(patsubst "%",%,$(lastword $(shell grep 'version =' gpl.go)))
RELEASE   = ./releases/$(VERSION)

GITHUB_USERNAME = "Code-Hex"

rm-build:
	@rm -rf build

rm-releases:
	@rm -rf releases

rm-all: rm-build rm-releases

release: all
	@mkdir -p $(RELEASE)
	@for DIR in $(shell ls ./build/$(VERSION)/) ; do \
		echo Processing in build/$(VERSION)/$$DIR; \
		cd $(PWD); \
		cp README.md ./build/$(VERSION)/$$DIR; \
		cp LICENSE ./build/$(VERSION)/$$DIR; \
		tar -cjf ./$(RELEASE)/gpl_$(VERSION)_$$DIR.tar.bz2 -C ./build/$(VERSION) $$DIR; \
		tar -czf ./$(RELEASE)/gpl_$(VERSION)_$$DIR.tar.gz -C ./build/$(VERSION) $$DIR; \
		cd build/$(VERSION); \
		zip -9 $(PWD)/$(RELEASE)/gpl_$(VERSION)_$$DIR.zip $$DIR/*; \
	done

prepare-github: github-token
	@echo "'github-token' file is required"

release-upload: prepare-github release
	@echo "Uploading..."
	@ghr -u $(GITHUB_USERNAME) -t $(shell cat github-token) --draft --replace $(VERSION) $(RELEASE)

all: test
	@PATH=$(SUB_BIN)/$(GOOS)/$(GOARCH):$(PATH)
	@gox -os=$(XC_OS) -arch=$(XC_ARCH) -output="build/$(VERSION)/{{.OS}}_{{.Arch}}/{{.Dir}}" $(CMDPATH)

test: deps
	@PATH=$(SUB_BIN)/$(GOOS)/$(GOARCH):$(PATH) go test -v $(shell glide nv)

deps: glide
	@PATH=$(SUB_BIN)/$(GOOS)/$(GOARCH):$(PATH) glide install
	go get github.com/golang/lint/golint
	go get github.com/mattn/goveralls
	go get github.com/axw/gocov/gocov
	go get github.com/mitchellh/gox

$(SUB_BIN)/$(GOOS)/$(GOARCH)/glide:
ifndef HAS_GLIDE
	@mkdir -p $(SUB_BIN)/$(GOOS)/$(GOARCH)
	@curl -L https://github.com/Masterminds/glide/releases/download/v0.11.0/glide-v0.11.0-$(GOOS)-$(GOARCH).zip -o glide.zip
	@unzip glide.zip
	@mv ./$(GOOS)-$(GOARCH)/glide $(SUB_BIN)/$(GOOS)/$(GOARCH)/glide
	@rm -rf ./$(GOOS)-$(GOARCH)
	@rm ./glide.zip
endif

glide: $(SUB_BIN)/$(GOOS)/$(GOARCH)/glide

lint: deps
	@for dir in $$(glide novendor); do \
	golint $$dir; \
	done;

cover: deps
	goveralls

.PHONY: test deps lint cover
