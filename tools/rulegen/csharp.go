package main

var csharpLibraryUsageTemplateString = `load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load(
    "@io_bazel_rules_dotnet//dotnet:defs.bzl",
    "core_register_sdk",
    "dotnet_register_toolchains",
    "dotnet_repositories",
)

core_version = "v2.1.503"

dotnet_register_toolchains(
    core_version = core_version,
)

dotnet_register_toolchains(
    core_version = core_version,
)

core_register_sdk(
    name = "core_sdk",
    core_version = core_version,
)

dotnet_repositories()

load("@build_stack_rules_proto//csharp/nuget:packages.bzl", nuget_packages = "packages")

nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")

nuget_protobuf_packages()`

var csharpProtoLibraryUsageTemplate = mustTemplate(csharpLibraryUsageTemplateString)

var csharpGrpcLibraryUsageTemplate = mustTemplate(csharpLibraryUsageTemplateString + `

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_grpc_packages")

nuget_grpc_packages()`)

var csharpLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )
`

var csharpProtoLibraryRuleTemplate = mustTemplate(csharpLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    core_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
        ],
        visibility = kwargs.get("visibility"),
    )`)

var csharpGrpcLibraryRuleTemplate = mustTemplate(csharpLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    core_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
            "@grpc.core//:netstandard1.5_core",
            "@system.interactive.async//:netstandard2.0_core",
        ],
        visibility = kwargs.get("visibility"),
    )`)

var csharpLibraryFlags = []*Flag{
	{
		Category:    "build",
		Name:        "strategy",
		Value:       "CoreCompile=standalone",
		Description: "dotnet SDK desperately wants to find the HOME directory",
	},
	{
		Category: "build",
		Name:     "incompatible_disallow_struct_provider_syntax",
		Value:    "false",
	},
}

func makeCsharp() *Language {
	return &Language{
		Dir:   "csharp",
		Name:  "csharp",
		Flags: commonLangFlags,
		Notes: mustTemplate(`**NOTE 1**: the csharp_* rules currently don't play nicely with sandboxing.  You may see errors like:

~~~python
The user's home directory could not be determined. Set the 'DOTNET_CLI_HOME' environment variable to specify the directory to use.
~~~

or

~~~python
System.ArgumentNullException: Value cannot be null.
Parameter name: path1
   at System.IO.Path.Combine(String path1, String path2)
   at Microsoft.DotNet.Configurer.CliFallbackFolderPathCalculator.get_DotnetUserProfileFolderPath()
   at Microsoft.DotNet.Configurer.FirstTimeUseNoticeSentinel..ctor(CliFallbackFolderPathCalculator cliFallbackFolderPathCalculator)
   at Microsoft.DotNet.Cli.Program.ProcessArgs(String[] args, ITelemetry telemetryClient)
   at Microsoft.DotNet.Cli.Program.Main(String[] args)
~~~

To remedy this, use --strategy=CoreCompile=standalone for the csharp rules (put it in your .bazelrc file).

**NOTE 2**: the csharp nuget dependency sha256 values do not appear stable.`),
		Rules: []*Rule{
			&Rule{
				Name:           "csharp_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//csharp:csharp"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates csharp protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "csharp_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//csharp:csharp", "//csharp:grpc_csharp"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates csharp protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "csharp_proto_library",
				Kind:           "proto",
				Implementation: csharpProtoLibraryRuleTemplate,
				Usage:          csharpProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates csharp protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          csharpLibraryFlags,
			},
			&Rule{
				Name:           "csharp_grpc_library",
				Kind:           "grpc",
				Implementation: csharpGrpcLibraryRuleTemplate,
				Usage:          csharpGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates csharp protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          csharpLibraryFlags,
			},
		},
	}
}
