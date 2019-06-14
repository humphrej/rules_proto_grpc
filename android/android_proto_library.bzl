load("//android:android_proto_compile.bzl", "android_proto_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def android_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    android_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create android library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            str(Label("//android:proto_deps")),
        ],
        exports = [
            str(Label("//android:proto_deps")),
        ],
        visibility = kwargs.get("visibility"),
    )
