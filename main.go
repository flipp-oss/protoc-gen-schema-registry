package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/flipp-oss/protoc-gen-schema-registry/input"
	"github.com/flipp-oss/protoc-gen-schema-registry/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	os.Setenv("GOLANG_PROTOBUF_REGISTRATION_CONFLICT", "warn")
	req, err := input.ReadRequest()
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("error reading request: %w", err))
	}
	params := input.ParseParams(req)
	registry.Setup(params)

	for _, file := range req.ProtoFile {
		if !slices.Contains(req.FileToGenerate, *file.Name) {
			continue
		}
		registry.Process(file, req.ProtoFile)
	}
	writeResponse()
}

func writeResponse() {
	feature := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	response := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &feature,
	}
	out, err := proto.Marshal(response)
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("error marshalling response: %w", err))
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("error writing response: %w", err))
	}
}
