SHELL=/bin/bash
.PHONY lint:
lint:
	golangci-lint run cocotola-api/src cocotola-translator-api/src

.PHONY gen-src:
gen-src:
	@pushd ./cocotola-api/ && \
	go generate ./src/... && \
	popd

# update-deps:
# 	bazel run //tools/update-deps
work-init:
	go work init cocotola-api cocotola-translator-api lib cocotola-synthesizer-api cocotola-tatoeba-api test-lib
gazelle:
	bazel run //:gazelle -- update-repos -from_file ./go.work
	bazel run //:gazelle

build:
	bazel build //...

run-api:
	bazel run //cocotola-api/src

test:
	bazel test //... --test_output=all

docker-build:
	bazel build //src:go_image

docker-run:
	bazel run //src:go_image

# work-init:
# 	go work cocotola-api cocotola-translator-api lib

test-docker-up:
	@docker-compose -f docker/test/docker-compose.yml up -d
	sleep 10
	@chmod -R 777 docker/test

test-docker-down:
	@docker-compose -f docker/test/docker-compose.yml down
