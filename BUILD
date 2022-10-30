load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/kujilabo/cocotola
gazelle(name = "gazelle")

go_library(
    name = "cocotola_lib",
    srcs = ["main.go"],
    importpath = "github.com/kujilabo/cocotola",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "cocotola",
    embed = [":cocotola_lib"],
    visibility = ["//visibility:public"],
)
