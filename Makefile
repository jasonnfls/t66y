GOPATH := $(shell pwd)
export GOPATH

.PHONY: clean all listing page dep

all: listing page

listing: bin/listing
bin/listing: src/listing.go src/common/utils.go
	go build -o bin/listing src/listing.go

page: bin/page
bin/page: src/page.go src/common/utils.go
	go build -o bin/page src/page.go

clean:
	rm -f bin/*

dep:
	mkdir -p bin
	go get -u "github.com/PuerkitoBio/goquery"
	go get -u "golang.org/x/text/encoding/simplifiedchinese"
	go get -u "golang.org/x/text/transform"
