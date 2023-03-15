# Go parameters
GOCMD:=$(shell which go)
GOLINT:=$(shell which golint)
GOIMPORT:=$(shell which goimports)
GOFMT:=$(shell which gofmt)
GOBUILD:=$(GOCMD) build
GORUN:=$(GOCMD) run
GOCLEAN:=$(GOCMD) clean
GOTEST:=$(GOCMD) test
GOGET:=$(GOCMD) get
GOLIST:=$(GOCMD) list
GOVET:=$(GOCMD) vet

# Generate HTML files
doc:
	$(GORUN) main.go -f=html

# Generate TypeScript files
ts:
	rm -rf ./output/ts/*
	$(GORUN) main.go -f=ts

# Generate CoreLibModule
buildNgLib:
	rm -rf ./ng-workspace/projects/ng-core-lib/src/lib/*
	mkdir -p ./ng-workspace/projects/ng-core-lib/src/lib
	cp -r ./output/ts/* ./ng-workspace/projects/ng-core-lib/src/lib
	npm run build-lib --prefix ./ng-workspace