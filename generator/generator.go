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

type Generator struct {
	Request  *plugin.CodeGeneratorRequest
	Response *plugin.CodeGeneratorResponse

	Verbose bool

	InputFileDescriptors []*descriptor.FileDescriptorProto

	PackageName      string
	PackageDirectory string

	MessageTypes MessageTypes

	ErlModuleName string
	ErlModulePath string

	erlModuleTemplate *template.Template
}

func NewGenerator(req *plugin.CodeGeneratorRequest) (*Generator, error) {
	g := Generator{
		Request:  req,
		Response: new(plugin.CodeGeneratorResponse),

		Verbose: true,
	}

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
	g.ErlModulePath = path.Join(g.PackageDirectory, g.ErlModuleName+".erl")

	err := g.generateFile(g.ErlModulePath, g.erlModuleTemplate, g)
	if err != nil {
		return fmt.Errorf("cannot generate erlang module: %w", err)
	}

	return nil
}

func (g *Generator) collectData() error {
	fns := []func() error{
		g.collectInputFileDescriptors,
		g.collectPackageName,
		g.collectPackageDirectory,
		g.collectMessageTypes,
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

	var addType func(*descriptor.FileDescriptorProto, *descriptor.DescriptorProto, *MessageType) error
	addType = func(fd *descriptor.FileDescriptorProto, d *descriptor.DescriptorProto, parent *MessageType) error {
		var mt MessageType
		if err := mt.FromDescriptor(fd, d, parent); err != nil {
			return fmt.Errorf("cannot create type for "+
				"message %s in package %s: %w",
				d.GetName(), fd.GetPackage(), err)
		}
		mts = append(mts, &mt)

		for _, nd := range d.NestedType {
			if err := addType(fd, nd, &mt); err != nil {
				return err
			}
		}

		return nil
	}

	for _, fd := range g.InputFileDescriptors {
		for _, d := range fd.MessageType {
			if err := addType(fd, d, nil); err != nil {
				return err
			}
		}
	}

	g.MessageTypes = mts
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
