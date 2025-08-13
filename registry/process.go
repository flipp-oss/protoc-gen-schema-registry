package registry

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"github.com/flipp-oss/protoc-gen-schema-registry/input"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

var protosToIgnore = []string{
	"google/protobuf/descriptor.proto",
}

func DirectNameStrategy(topic string, serdeType serde.Type, schema schemaregistry.SchemaInfo) (string, error) {
	return topic, nil
}

var client schemaregistry.Client
var ser *protobuf.Serializer

func Setup(params input.Params) {
	var err error
	client, err = schemaregistry.NewClient(schemaregistry.NewConfig(params.RegistryUrl))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Created schema registry client for %s\n", params.RegistryUrl)

	ser, err = protobuf.NewSerializer(client, serde.ValueSerde, protobuf.NewSerializerConfig())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create serializer: %s\n", err)
		os.Exit(1)
	}
	ser.SubjectNameStrategy = DirectNameStrategy

	fmt.Fprintf(os.Stderr, "Created protobuf serializer\n")
}

var processedFiles = map[string]bool{}

func Process(fileProto *descriptorpb.FileDescriptorProto, allFiles []*descriptorpb.FileDescriptorProto) {
	packagePath := strings.Replace(fileProto.GetPackage(), ".", "/", -1)
	subject := fileProto.GetName()
	if !strings.HasPrefix(subject, packagePath) {
		subject = subject + "/" + fileProto.GetName()
	}
	if slices.Contains(protosToIgnore, subject) {
		return
	}
	if processedFiles[subject] {
		return
	}
	processedFiles[subject] = true

	for _, dependency := range fileProto.GetDependency() {
		for _, file := range allFiles {
			if file.GetName() == dependency {
				Process(file, allFiles)
			}
		}
	}
	fileDesc, err := protodesc.NewFile(fileProto, protoregistry.GlobalFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file descriptor for %s: %s\n", fileProto.GetName(), err)
	}
	fmt.Fprintf(os.Stderr, "Created file descriptor for %s\n", subject)
	_, err = protoregistry.GlobalFiles.FindFileByPath(subject)
	if err != nil {
		err = protoregistry.GlobalFiles.RegisterFile(fileDesc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to register file descriptor for %s: %s\n", fileProto.GetName(), err)
		}
	}
	if fileDesc.Messages().Len() > 0 {
		message := dynamicpb.NewMessage(fileDesc.Messages().Get(0))
		_, err = ser.Serialize(subject, message)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to serialize message: %s\n", err)
		}
	}
}
