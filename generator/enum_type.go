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
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type EnumType struct {
	Parent *MessageType

	Package  string
	Name     string
	FullName string

	ErlName string
}

type EnumTypes []*EnumType

func (enumType *EnumType) FromDescriptor(fd *descriptor.FileDescriptorProto, ed *descriptor.EnumDescriptorProto, parent *MessageType) error {
	et := EnumType{
		Parent: parent,

		Package: fd.GetPackage(),
		Name:    ed.GetName(),
	}

	et.FullName = EnumTypeFullName(&et)
	et.ErlName = EnumTypeFullNameToErlTypeName(et.FullName)

	*enumType = et
	return nil
}

func EnumTypeFullName(et *EnumType) string {
	if et.Parent == nil {
		return et.Name
	}

	return MessageTypeFullName(et.Parent) + "." + et.Name
}

func EnumTypeFullNameToErlTypeName(name string) string {
	name2 := strings.ReplaceAll(name, ".", "__")
	return CamelCaseToSnakeCase(name2)
}
