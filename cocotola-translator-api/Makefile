.PHONY: test
test:
	@go test -coverprofile="coverage.txt" -covermode=atomic ./...

lint:
	# @golangci-lint run --disable-all --enable wrapcheck       --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable gocognit        --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable gocyclo         --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable errorlint       --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable gomnd           --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable unconvert       --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable gofmt           --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable gosec           --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable noctx           --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable errcheck        --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable whitespace      --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable forcetypeassert --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable forbidigo       --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable bodyclose       --config ../.github/.golangci.yml --skip-files ".*_test.go"
	@golangci-lint run --disable-all --enable exhaustive      --config ../.github/.golangci.yml --skip-files ".*_test.go"
