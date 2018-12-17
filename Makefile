
binary	:= server
src=$(shell find . -type f -name '*.go')
tests=$(shell find . -type f -name '*_test.go')

.PHONY: help lint clean

build: $(binary) client/sws.min.js  ## Build a binary

help: ## This help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |sort |awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

$(binary): $(src)
	go build -ldflags "-w -s -X main.version=$(version)" \
		-o $(binary) ./cmd/$(binary)

%.min.js: %.js
	yarn run uglifyjs -o $@ $<

test: $(tests) $(src) ## Run tests
	go test -v -short -coverprofile=coverage.out -cover ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@for file in $$(find . -name 'vendor' -prune -o -type f -name '*.go'); do golint $$file; done

clean: ## Clean all test files
	rm -f $(binary)
	rm -rf coverage*
