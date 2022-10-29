.PHONY lint:
lint:
	golangci-lint run cocotola-api cocotola-translator-api

update-deps:
	bazel run //tools/update-deps

gazelle:
	bazel run //:gazelle -- update-repos -from_file ./go.mod
	bazel run //:gazelle

build:
	bazel build //...

run:
	bazel run //...

run-api:
	bazel run //cocotola-api
run-translator-api:
	bazel run //cocotola-translator-api

test:
	bazel test //... --test_output=all

docker-build:
	bazel build //src:go_image

docker-run:
	bazel run //src:go_image

work-init:
	go work cocotola-api cocotola-translator-api lib
