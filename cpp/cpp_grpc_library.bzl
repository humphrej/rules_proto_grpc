load("//cpp:cpp_grpc_compile.bzl", "cpp_grpc_compile")

def cpp_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    cpp_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create cpp library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
            "@com_github_grpc_grpc//:grpc++",
            "@com_github_grpc_grpc//:grpc++_reflection",
        ],
        # This seems magical to me.
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )
