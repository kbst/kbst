test:
	go test --cover ./...

build: test
	go build -o kbst .

install: test
	go install .

snapshot: test
	goreleaser release --skip-publish --snapshot --clean --skip-sign

tidy:
	go mod tidy
