.PHONY: all gox cli cli-linux-amd64 cli-linux-arm64 cli-darwin-amd64 cli-darwin-arm64 build format lint tidy test-setup test integration-test clean

# -----------------------------------------------------------------------------
#  CONSTANTS
# -----------------------------------------------------------------------------

version = `git describe --abbrev=0 --tags | tr -d 'v'`

go_dirs = $(shell go list ./... | grep -v /build/ | grep -v /vendor/)

build_dir       = build
integration_dir = integration

coverage_dir  = $(build_dir)/coverage
coverage_out  = $(coverage_dir)/coverage.out
coverage_html = $(coverage_dir)/coverage.html

output_dir   = $(build_dir)/output
compiled_dir = $(build_dir)/compiled

linux_amd64_out_dir   = $(output_dir)/linux/amd64
linux_arm64_out_dir   = $(output_dir)/linux/arm64
darwin_amd64_out_dir  = $(output_dir)/darwin/amd64
darwin_arm64_out_dir  = $(output_dir)/darwin/arm64

bin_name          = conntest
bin_linux_amd64   = $(linux_amd64_out_dir)/$(bin_name)
bin_linux_arm64   = $(linux_arm64_out_dir)/$(bin_name)
bin_darwin_amd64  = $(darwin_amd64_out_dir)/$(bin_name)
bin_darwin_arm64  = $(darwin_arm64_out_dir)/$(bin_name)

# -----------------------------------------------------------------------------
#  BUILDING
# -----------------------------------------------------------------------------

all: cli

gox:
	go install github.com/mitchellh/gox@latest
	mkdir -p $(compiled_dir)

cli: gox cli-linux-amd64 cli-linux-arm64 cli-darwin-amd64 cli-darwin-arm64
	(cd $(linux_amd64_out_dir) && zip -r staging.zip $(bin_name))
	mv $(linux_amd64_out_dir)/staging.zip $(compiled_dir)/conntest_$(version)_linux_amd64.zip
	(cd $(linux_arm64_out_dir) && zip -r staging.zip $(bin_name))
	mv $(linux_arm64_out_dir)/staging.zip $(compiled_dir)/conntest_$(version)_linux_arm64.zip
	(cd $(darwin_amd64_out_dir) && zip -r staging.zip $(bin_name))
	mv $(darwin_amd64_out_dir)/staging.zip $(compiled_dir)/conntest_$(version)_darwin_amd64.zip
	(cd $(darwin_arm64_out_dir) && zip -r staging.zip $(bin_name))
	mv $(darwin_arm64_out_dir)/staging.zip $(compiled_dir)/conntest_$(version)_darwin_arm64.zip

cli-linux-amd64: gox
	CGO_ENABLED=0 gox -osarch=linux/amd64 -output=$(bin_linux_amd64) -ldflags="-X main.version=$(version)" .

cli-linux-arm64: gox
	CGO_ENABLED=0 gox -osarch=linux/arm64 -output=$(bin_linux_arm64) -ldflags="-X main.version=$(version)" .

cli-darwin-amd64: gox
	CGO_ENABLED=0 gox -osarch=darwin/amd64 -output=$(bin_darwin_amd64) -ldflags="-X main.version=$(version)" .

cli-darwin-arm64: gox
	CGO_ENABLED=0 gox -osarch=darwin/arm64 -output=$(bin_darwin_arm64) -ldflags="-X main.version=$(version)" .

# Simple build for local development
build:
	go build -ldflags="-X main.version=$(version)" -o $(bin_name) .

# -----------------------------------------------------------------------------
#  FORMATTING
# -----------------------------------------------------------------------------

format:
	go fmt ./...
	gofmt -s -w .

lint:
	go install golang.org/x/lint/golint@latest
	golint ./...

tidy:
	go mod tidy

update:
	go get -u ./...

# -----------------------------------------------------------------------------
#  TESTING
# -----------------------------------------------------------------------------

test-setup:
	mkdir -p $(coverage_dir)
	go install golang.org/x/tools/cmd/cover@latest

test: test-setup
	go test $(go_dirs) -v -short -covermode=count -coverprofile=$(coverage_out)
	go tool cover -html=$(coverage_out) -o $(coverage_html)
	go tool cover -func=$(coverage_out)

integration-test: test-setup
	go test $(go_dirs) -v -covermode=count -coverprofile=$(coverage_out)
	go tool cover -html=$(coverage_out) -o $(coverage_html)
	go tool cover -func=$(coverage_out)

# -----------------------------------------------------------------------------
#  CLEANUP
# -----------------------------------------------------------------------------

clean:
	rm -rf $(build_dir)
	rm -f $(bin_name)
