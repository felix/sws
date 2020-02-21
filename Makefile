
VERSION?=	$(shell git describe --tags --always)
BINARY=	$(patsubst %,dist/%,$(shell find cmd/* -maxdepth 0 -type d -exec basename {} \;))
SRC=	$(shell find . -type f -name '*.go')
SQL=	$(shell find sql -type f)
EXTRAS=	cmd/server/migrations.go \
	counter/sws.min.js \
	cmd/server/counter.go \
	cmd/server/static.go \
	cmd/server/templates.go
TMPL=	$(shell find templates -type f -name '*.tmpl')
STATIC=	$(shell find static -type f)

.PHONY: build
build: $(BINARY)

dist/%: $(SRC) $(EXTRAS)
	go build -ldflags "-X main.Version=$(VERSION)" -o $@ ./cmd/$*

cmd/server/static.go: $(STATIC)
	go generate ./static >$@

cmd/server/templates.go: $(TMPL)
	go generate ./templates >$@

cmd/server/counter.go: counter/sws.min.js
	@printf "package main\n\nconst counter = \`" >$@
	@cat $< >>$@
	@printf "\`\n" >>$@

# cmd/server/counter.go: counter/sws.min.js
# 	go generate ./counter >$@


cmd/server/migrations.go: $(SQL)
	go generate ./sql >$@

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
