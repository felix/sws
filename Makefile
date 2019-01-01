
binary	:= collector server
src		:= $(shell find . -type f -name '*.go')
extras	:= cmd/collector/snippet.go
tests	:= $(shell find . -type f -name '*_test.go')

.PHONY: help build lint clean

build: $(binary) ## Build a binary

$(binary): $(src) $(extras)
	go build -ldflags "-w -s -X main.version=$(version)" -o $@ ./cmd/$@

cmd/collector/snippet.go: client/sws.min.js
	@printf "package main\n\nconst snippet = \`" >$@
	@cat $< >>$@
	@printf "\`\n" >>$@

%.min.js: %.js node_modules
	@yarn run uglifyjs -o $@ $<

node_modules: package.json
	yarn

test: lint $(tests) $(src) ## Run tests
	go test -v -short -coverprofile=coverage.out -cover ./...
	go tool cover -html=coverage.out -o coverage.html

lint: $(src)
	golint $<

clean: ## Clean all test files
	rm -f $(binary) $(extras)
	rm -f client/sws.min.js
	rm -rf coverage*

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort |awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
