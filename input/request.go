package input

import (
	"io"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

type Params struct {
	RegistryUrl string
}

func ReadRequest() (*pluginpb.CodeGeneratorRequest, error) {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	req := &pluginpb.CodeGeneratorRequest{}
	err = proto.Unmarshal(in, req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func parseRawParams(req *pluginpb.CodeGeneratorRequest) map[string]string {
	param := req.GetParameter()
	if len(param) == 0 {
		return nil
	}
	paramTokens := strings.Split(param, ",")
	paramMap := map[string]string{}
	for _, token := range paramTokens {
		paramStrings := strings.Split(token, "=")
		if len(paramStrings) == 2 {
			paramMap[paramStrings[0]] = paramStrings[1]
		}
	}
	return paramMap
}

func ParseParams(req *pluginpb.CodeGeneratorRequest) Params {
	params := Params{}
	rawParams := parseRawParams(req)
	for k, v := range rawParams {
		if k == "registry_url" {
			params.RegistryUrl = v
		}
	}
	return params
}
