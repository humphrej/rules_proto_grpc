load(
    "//:repositories.bzl",
    "io_bazel_rules_d",
    "com_github_dcarp_protobuf_d",
    "rules_proto_grpc_repos",
)

def d_repos(**kwargs):
    rules_proto_grpc_repos(**kwargs)
    com_github_dcarp_protobuf_d(**kwargs)
    io_bazel_rules_d(**kwargs)
