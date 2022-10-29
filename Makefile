gazelle:
	bazel run //:gazelle -- update-repos -from_file ./go.mod
	bazel run //:gazelle

build:
	bazel build //src

run:
	bazel run //src

test:
	bazel test //... --test_output=all

docker-build:
	bazel build //src:go_image

docker-run:
	bazel run //src:go_image
