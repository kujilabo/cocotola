SHELL=/bin/bash

.PHONY: all
all:
	# $(MAKE) lint
	$(MAKE) gen-swagger
	$(MAKE) gen-src
	$(MAKE) gen-proto
	$(MAKE) gazelle
	$(MAKE) update-mod
	$(MAKE) build
	$(MAKE) test
	$(MAKE) dev-docker-build

.PHONY: lint
lint:
	@pushd ./cocotola-api/ && \
		docker run --rm -i hadolint/hadolint < Dockerfile && \
	popd
	@pushd ./cocotola-synthesizer-api/ && \
		docker run --rm -i hadolint/hadolint < Dockerfile && \
	popd
	@pushd ./cocotola-tatoeba-api/ && \
		docker run --rm -i hadolint/hadolint < Dockerfile && \
	popd
	@pushd ./cocotola-translator-api/ && \
		docker run --rm -i hadolint/hadolint < Dockerfile && \
	popd

	@pushd ./cocotola-api/src && \
		golangci-lint run --config ../../.github/.golangci.yml && \
	popd
	@pushd ./cocotola-synthesizer-api/src && \
		golangci-lint run --config ../../.github/.golangci.yml && \
	popd
	@pushd ./cocotola-tatoeba-api/src && \
		golangci-lint run --config ../../.github/.golangci.yml && \
	popd
	@pushd ./cocotola-translator-api/src && \
		golangci-lint run --config ../../.github/.golangci.yml && \
	popd

.PHONY: gen-swagger
gen-swagger:
	@pushd ./cocotola-api/ && \
		swag init -d src -o src/docs && \
	popd
	@pushd ./cocotola-synthesizer-api/ && \
		swag init -d src -o src/docs && \
	popd
	@pushd ./cocotola-tatoeba-api/ && \
		swag init -d src -o src/docs && \
	popd
	@pushd ./cocotola-translator-api/ && \
		swag init -d src -o src/docs && \
	popd

.PHONY: gen-src
gen-src:
	@pushd ./cocotola-api/ && \
		go generate ./src/... && \
	popd
	@pushd ./cocotola-synthesizer-api/ && \
		go generate ./src/... && \
	popd
	@pushd ./cocotola-tatoeba-api/ && \
		go generate ./src/... && \
	popd
	@pushd ./cocotola-translator-api/ && \
		go generate ./src/... && \
	popd

.PHONY: gen-proto
gen-proto:
	@pushd ./proto && \
	protoc --go_out=./src/ --go_opt=paths=source_relative \
        --go-grpc_out=./src/ --go-grpc_opt=paths=source_relative \
        proto/translator_admin.proto \
		proto/translator_user.proto && \
	popd

.PHONY: update-mod
update-mod:
	@pushd ./cocotola-api/ && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-synthesizer-api/ && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-tatoeba-api/ && \
		go get -u ./... && \
	popd
	@pushd ./cocotola-translator-api/ && \
		go get -u ./... && \
	popd

# update-deps:
# 	bazel run //tools/update-deps
work-init:
	@go work init
	@go work use -r .

gazelle:
	sudo chmod 777 -R docker/development
	@bazel run //:gazelle -- update-repos -from_file ./go.work
	@bazel run //:gazelle

build:
	@bazel build //...

run-cocotola:
	@bazel run //cocotola-api/src

test:
	@bazel test //... --test_output=all --test_timeout=60

docker-build:
	bazel build //src:go_image

docker-run:
	bazel run //src:go_image

# work-init:
# 	go work cocotola-api cocotola-translator-api lib

dev-docker-up:
	@docker-compose -f docker/development/docker-compose.yml up -d
	sleep 10
	@chmod -R 777 docker/test

dev-docker-down:
	@docker-compose -f docker/development/docker-compose.yml down
	sleep 10
	@chmod -R 777 docker/test

dev-docker-clean:
	@rm -rf docker/development/mysql-*

test-docker-up:
	@docker-compose -f docker/test/docker-compose.yml up -d
	sleep 10
	@chmod -R 777 docker/test

test-docker-down:
	@docker-compose -f docker/test/docker-compose.yml down

.PHONY: dev-docker-build
dev-docker-build:
	@pushd ./cocotola-api/ && \
		docker build -t cocotola-api . && \
	popd
	@pushd ./cocotola-synthesizer-api/ && \
		docker build -t cocotola-synthesizer-api . && \
	popd
	@pushd ./cocotola-tatoeba-api/ && \
		docker build -t cocotola-tatoeba-api . && \
	popd
	@pushd ./cocotola-translator-api/ && \
		docker build -t cocotola-translator-api . && \
	popd
