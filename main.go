package main

import (
	"flag"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.4"

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-simple %v\n", version)
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			//generateFile(gen, f)
			if len(f.Messages) > 0 {
				for _, message := range f.Messages {
					if _, found := strings.CutSuffix(string(message.Desc.Name()), "Model"); found {
						generateModelFile(gen, f, message)
						generateTableFile(gen, f, message)
					}
				}
			}
			if len(f.Services) > 0 {
				for _, service := range f.Services {
					generateSimpleServerCode(gen, f, service)
					generateApiCode(gen, f, service)
				}
			}

		}
		return nil
	})
}
