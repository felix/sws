
VERSION?=	$(shell git describe --tags --always)
BINARY=	$(patsubst %,dist/%,$(shell find cmd/* -maxdepth 0 -type d -exec basename {} \;))
SRC=	$(shell find . -type f -name '*.go')
SQL=	$(shell find sql -type f)
EXTRAS=	cmd/server/migrations.go \
	counter/sws.min.js \
	cmd/server/counter.go \
	cmd/server/templates.go
TMPL=	$(shell find tmpl -type f -name '*.tmpl')
STATIC=	static/default.css

.PHONY: build
build: $(BINARY)

dist/%: $(SRC) $(EXTRAS)
	go build -ldflags "-X main.Version=$(VERSION)" -o $@ ./cmd/$*

cmd/server/templates.go: $(TMPL)  $(STATIC)
	go generate ./scripts >$@

cmd/server/migrations.go: $(SQL)
	go generate ./sql >$@

cmd/server/counter.go: counter/sws.min.js
	printf "package main\n\nfunc getCounter() string { return \`" >$@
	cat $< >>$@
	printf "\`}\n" >>$@

static/default.css: sass/main.scss
	yarn run -s node-sass $< $@

%.min.js: %.js node_modules
	yarn run -s uglifyjs -c -m -o $@ $<

node_modules: package.json
	yarn -s

.PHONY: test
test: lint
	go test -v -short -coverprofile=coverage.out -cover ./... \
		&& go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint:
	go vet ./...

.PHONY: clean
clean:
	rm -fr dist
	rm -f client/sws.min.js
	rm -f $(EXTRAS)
	rm -fr coverage*

.PHONY: dist-clean
dist-clean: clean
	rm -fr node_modules
