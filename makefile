clean:
	rm sw-crawl

build:
	@echo building linux
	env CGO_ENABLED=0 go build

run: build
	./sw-crawl