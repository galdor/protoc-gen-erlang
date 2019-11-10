// Copyright (c) 2019 Nicolas Martyanoff <khaelin@gmail.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type AbsoluteNameResolver interface {
	FindMessageType(string) *MessageType
	FindEnumType(string) *EnumType
}

type Generator struct {
	Request  *plugin.CodeGeneratorRequest
	Response *plugin.CodeGeneratorResponse

	Verbose bool

	InputFileDescriptors []*descriptor.FileDescriptorProto

	PackageName      string
	PackageDirectory string

	MessageTypes              MessageTypes
	PackageMessageTypes       MessageTypes
	DescriptorToMessageType   map[*descriptor.DescriptorProto]*MessageType
	AbsoluteNameToMessageType map[string]*MessageType

	EnumTypes              EnumTypes
	PackageEnumTypes       EnumTypes
	AbsoluteNameToEnumType map[string]*EnumType

	ErlModuleName string
	ErlHRLPath    string
	ErlModulePath string

	erlHRLTemplate    *template.Template
	erlModuleTemplate *template.Template
}

func NewGenerator(req *plugin.CodeGeneratorRequest) (*Generator, error) {
	g := Generator{
		Request:  req,
		Response: new(plugin.CodeGeneratorResponse),

		Verbose: true,
	}

	erlHRLTemplate, err := ErlHRLTemplate()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot create erlang hrl template: %w", err)
	}
	g.erlHRLTemplate = erlHRLTemplate

	erlModuleTemplate, err := ErlModuleTemplate()
	if err != nil {
		return nil, fmt.Errorf(
			"cannot create erlang module template: %w", err)
	}
	g.erlModuleTemplate = erlModuleTemplate

	return &g, nil
}

func (g *Generator) Info(format string, args ...interface{}) {
	if !g.Verbose {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (g *Generator) GenerateOutput() error {
	if err := g.collectData(); err != nil {
		return err
	}

	g.ErlModuleName = ProtoPackageNameToErlModuleName(g.PackageName)

	g.ErlHRLPath = path.Join(g.PackageDirectory, g.ErlModuleName+".hrl")
	g.ErlModulePath = path.Join(g.PackageDirectory, g.ErlModuleName+".erl")

	var err error

	err = g.generateFile(g.ErlHRLPath, g.erlHRLTemplate, g)
	if err != nil {
		return fmt.Errorf("cannot generate erlang hrl file: %w", err)
	}

	err = g.generateFile(g.ErlModulePath, g.erlModuleTemplate, g)
	if err != nil {
		return fmt.Errorf("cannot generate erlang module: %w", err)
	}

	return nil
}

func (g *Generator) FindMessageType(absName string) *MessageType {
	mt, found := g.AbsoluteNameToMessageType[absName]
	if !found {
		return nil
	}

	return mt
}

func (g *Generator) FindEnumType(absName string) *EnumType {
	et, found := g.AbsoluteNameToEnumType[absName]
	if !found {
		return nil
	}

	return et
}

func (g *Generator) collectData() error {
	fns := []func() error{
		g.collectInputFileDescriptors,
		g.collectPackageName,
		g.collectPackageDirectory,
		g.collectMessageTypes,
		g.collectEnumTypes,
		g.resolveTypes,
	}

	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) collectInputFileDescriptors() error {
	var fds []*descriptor.FileDescriptorProto

	for _, name := range g.Request.FileToGenerate {
		for _, fd := range g.Request.ProtoFile {
			if fd.GetName() == name {
				fds = append(fds, fd)
			}
		}
	}

	g.InputFileDescriptors = fds

	return nil
}

func (g *Generator) collectPackageName() error {
	var name string

	for _, fd := range g.InputFileDescriptors {
		if name != "" && name != fd.GetPackage() {
			return errors.New("cannot process multiple files " +
				"with different proto packages")
		}

		name = fd.GetPackage()
	}

	g.PackageName = name
	return nil
}

func (g *Generator) collectPackageDirectory() error {
	var dir string

	for _, fd := range g.InputFileDescriptors {
		fdDir := path.Dir(fd.GetName())
		if dir != "" && dir != fdDir {
			return errors.New("cannot process multiple files " +
				"from different directories")
		}

		dir = fdDir
	}

	g.PackageDirectory = dir
	return nil
}

func (g *Generator) collectMessageTypes() error {
	var mts MessageTypes

	descriptorToMessageType := make(map[*descriptor.DescriptorProto]*MessageType)
	absoluteNameToMessageType := make(map[string]*MessageType)

	var addType func(*descriptor.FileDescriptorProto, *descriptor.DescriptorProto, *MessageType) error
	addType = func(fd *descriptor.FileDescriptorProto, d *descriptor.DescriptorProto, parent *MessageType) error {
		var mt MessageType
		if err := mt.FromDescriptor(fd, d, parent); err != nil {
			return fmt.Errorf("cannot create type for "+
				"message %s in package %s: %w",
				d.GetName(), fd.GetPackage(), err)
		}

		mts = append(mts, &mt)

		descriptorToMessageType[d] = &mt
		absoluteNameToMessageType[mt.AbsoluteName] = &mt

		for _, nd := range d.NestedType {
			if err := addType(fd, nd, &mt); err != nil {
				return err
			}
		}

		return nil
	}

	for _, fd := range g.Request.ProtoFile {
		for _, d := range fd.MessageType {
			if err := addType(fd, d, nil); err != nil {
				return err
			}
		}
	}

	g.MessageTypes = mts
	g.DescriptorToMessageType = descriptorToMessageType
	g.AbsoluteNameToMessageType = absoluteNameToMessageType

	for _, mt := range g.MessageTypes {
		if mt.Package == g.PackageName {
			g.PackageMessageTypes = append(g.PackageMessageTypes, mt)
		}
	}

	return nil
}

func (g *Generator) collectEnumTypes() error {
	var ets EnumTypes

	absoluteNameToEnumType := make(map[string]*EnumType)

	var addType func(*descriptor.FileDescriptorProto, *descriptor.EnumDescriptorProto, *MessageType) error
	addType = func(fd *descriptor.FileDescriptorProto, ed *descriptor.EnumDescriptorProto, parent *MessageType) error {
		var et EnumType
		if err := et.FromDescriptor(fd, ed, parent); err != nil {
			return fmt.Errorf("cannot create type for "+
				"enum %s in package %s: %w",
				ed.GetName(), fd.GetPackage(), err)
		}

		ets = append(ets, &et)

		absoluteNameToEnumType[et.AbsoluteName] = &et

		return nil
	}

	var addNestedType func(*descriptor.FileDescriptorProto, *descriptor.DescriptorProto) error
	addNestedType = func(fd *descriptor.FileDescriptorProto, d *descriptor.DescriptorProto) error {
		mt, found := g.DescriptorToMessageType[d]
		if !found {
			return fmt.Errorf("no message type found for "+
				"message %s in package %s",
				d.GetName(), fd.GetPackage())
		}

		for _, ed := range d.EnumType {
			if err := addType(fd, ed, mt); err != nil {
				return err
			}
		}

		for _, nd := range d.NestedType {
			if err := addNestedType(fd, nd); err != nil {
				return err
			}
		}

		return nil
	}

	for _, fd := range g.Request.ProtoFile {
		for _, ed := range fd.EnumType {
			if err := addType(fd, ed, nil); err != nil {
				return err
			}
		}

		for _, d := range fd.MessageType {
			if err := addNestedType(fd, d); err != nil {
				return err
			}
		}
	}

	g.EnumTypes = ets

	for _, et := range g.EnumTypes {
		if et.Package == g.PackageName {
			g.PackageEnumTypes = append(g.PackageEnumTypes, et)
		}
	}

	g.AbsoluteNameToEnumType = absoluteNameToEnumType

	return nil
}

func (g *Generator) resolveTypes() error {
	for _, mt := range g.MessageTypes {
		if err := mt.ResolveTypes(g); err != nil {
			return fmt.Errorf("cannot resolve types in message %q "+
				"of package %q: %w", mt.Name, mt.Package, err)
		}
	}

	return nil
}

func (g *Generator) generateFile(fileName string, tpl *template.Template, data interface{}) error {
	g.Info("generating %s", fileName)

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}

	content := buf.String()
	file := plugin.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &content,
	}

	g.Response.File = append(g.Response.File, &file)
	return nil
}
