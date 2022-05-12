//  Copyright (C) Isovalent, Inc. - All Rights Reserved.
//
//  NOTICE: All information contained herein is, and remains the property of
//  Isovalent Inc and its suppliers, if any. The intellectual and technical
//  concepts contained herein are proprietary to Isovalent Inc and its suppliers
//  and may be covered by U.S. and Foreign Patents, patents in process, and are
//  protected by trade secret or copyright law.  Dissemination of this information
//  or reproduction of this material is strictly forbidden unless prior written
//  permission is obtained from Isovalent Inc.
//

package common

import (
	"fmt"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewGeneratedFile creates a new codegen pakage and file in the project
func NewGeneratedFile(gen *protogen.Plugin, file *protogen.File, pkg string) *protogen.GeneratedFile {
	importPath := filepath.Join(string(file.GoImportPath), "codegen", pkg)
	fileName := filepath.Join(strings.TrimSuffix(file.GeneratedFilenamePrefix, "fgs"), "codegen", pkg, fmt.Sprintf("%s.pb.go", pkg))

	g := gen.NewGeneratedFile(fileName, protogen.GoImportPath(importPath))
	g.P(`//  Copyright (C) Isovalent, Inc. - All Rights Reserved.
//
//  NOTICE: All information contained herein is, and remains the property of
//  Isovalent Inc and its suppliers, if any. The intellectual and technical
//  concepts contained herein are proprietary to Isovalent Inc and its suppliers
//  and may be covered by U.S. and Foreign Patents, patents in process, and are
//  protected by trade secret or copyright law.  Dissemination of this information
//  or reproduction of this material is strictly forbidden unless prior written
//  permission is obtained from Isovalent Inc.
//
`)
	g.P("// Code generated by protoc-gen-go-tetragon. DO NOT EDIT")
	g.P()
	g.P("package ", pkg)
	g.P()

	return g
}

// GoIdent is a convenience helper that returns a qualified go ident as a string for
// a given import package and name
func GoIdent(g *protogen.GeneratedFile, importPath string, name string) string {
	return g.QualifiedGoIdent(protogen.GoIdent{
		GoName:       name,
		GoImportPath: protogen.GoImportPath(importPath),
	})
}

// FgsApiIdent is a convenience helper that calls GoIdent with the path to the FGS API
// package.
func FgsApiIdent(g *protogen.GeneratedFile, name string) string {
	return GoIdent(g, "github.com/isovalent/tetragon-oss/api/v1/fgs", name)
}

// Logger is a convenience helper that generates a call to logger.GetLogger()
func Logger(g *protogen.GeneratedFile) string {
	return fmt.Sprintf("%s()", GoIdent(g, "github.com/isovalent/tetragon-oss/pkg/logger", "GetLogger"))
}

// FmtErrorf is a convenience helper that generates a call to fmt.Errorf
func FmtErrorf(g *protogen.GeneratedFile, fmt_ string, args ...string) string {
	args = append([]string{fmt.Sprintf("\"%s\"", fmt_)}, args...)
	return fmt.Sprintf("%s(%s)", GoIdent(g, "fmt", "Errorf"), strings.Join(args, ", "))
}

// GetEvents returns a list of all messages that are events
func GetEvents(file *protogen.File) ([]*protogen.Message, error) {
	var getEventsResponse *protogen.Message
	for _, msg := range file.Messages {
		if msg.GoIdent.GoName == "GetEventsResponse" {
			getEventsResponse = msg
			break
		}
	}
	if getEventsResponse == nil {
		return nil, fmt.Errorf("Unable to find GetEventsResponse message")
	}

	var eventOneof *protogen.Oneof
	for _, oneof := range getEventsResponse.Oneofs {
		if oneof.Desc.Name() == "event" {
			eventOneof = oneof
			break
		}
	}
	if eventOneof == nil {
		return nil, fmt.Errorf("Unable to find GetEventsResponse.event")
	}

	validNames := make(map[string]struct{})
	for _, type_ := range eventOneof.Fields {
		name := strings.TrimPrefix(type_.GoIdent.GoName, "GetEventsResponse_")
		validNames[name] = struct{}{}
	}

	var events []*protogen.Message
	for _, msg := range file.Messages {
		if _, ok := validNames[string(msg.Desc.Name())]; ok {
			events = append(events, msg)
		}
	}

	return events, nil
}

// EventFieldCheck returns true if the event has the field
func EventFieldCheck(msg *protogen.Message, field string) bool {
	if msg.Desc.Fields().ByName(protoreflect.Name(field)) != nil {
		return true
	}

	return false
}

// IsProcessEvent returns true if the message is an FGS event that has a process field
func IsProcessEvent(msg *protogen.Message) bool {
	return EventFieldCheck(msg, "process")
}

// IsParentEvent returns true if the message is an FGS event that has a parent field
func IsParentEvent(msg *protogen.Message) bool {
	return EventFieldCheck(msg, "parent")
}
