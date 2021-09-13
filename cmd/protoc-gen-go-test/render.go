package main

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
)

var importPackages = map[string]string{
	"net":                                   "",
	"testing":                               "",
	"context":                               "",
	"github.com/stretchr/testify/assert":    "",
	"google.golang.org/grpc/test/bufconn":   "",
	"google.golang.org/grpc":                "",
	"google.golang.org/protobuf/proto":      "",
	"github.com/gotomicro/ego/client/egrpc": "cegrpc",
	"github.com/gotomicro/ego/server/egrpc": "",
}

// generateFile generates a _errors.pb.go file containing ego errors definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File, pkg string) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + "_test.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-go-test. DO NOT EDIT.")
	g.P()
	g.P("package ", pkg)
	g.P()
	g.P("import (")
	for path, name := range importPackages {
		g.P(fmt.Sprintf("%s \"%s\"", name, path))
	}
	g.P(fmt.Sprintf("%s %s", file.GoPackageName, file.GoImportPath))
	g.P(")")
	generateFileContent(gen, file, g)
	return g
}

// generateFileContent generates the ego errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the ego package it is being compiled against.")
	g.P()
	index := 0
	for _, svc := range file.Services {
		if generateTestSection(gen, file, g, svc) == false {
			index++
		}
	}
	// If all enums do not contain 'service', the current file is skipped
	if index == 0 {
		g.Skip()
	}
}

func generateTestSection(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, svc *protogen.Service) bool {
	var w svcWrapper
	for _, method := range svc.Methods {
		method.Desc.IsStreamingServer()
		d := &svcData{
			Name:              method.GoName,
			InType:            string(method.Input.Desc.Name()),
			OutType:           string(method.Output.Desc.Name()),
			Package:           string(file.GoPackageName),
			Service:           service{Name: svc.GoName},
			isStreamingServer: method.Desc.IsStreamingServer(),
			isStreamingClient: method.Desc.IsStreamingClient(),
		}
		w.Svcs = append(w.Svcs, d)
	}
	if len(w.Svcs) == 0 {
		return true
	}
	g.P(w.execute())
	return false
}
