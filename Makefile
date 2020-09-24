build:
	go build -o kbst .

install:
	go install .

snapshot:
	goreleaser release --skip-publish --snapshot --rm-dist --skip-sign

tidy:
	go mod tidy
