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

type MessageType struct {
	Parent *MessageType

	Package  string
	Name     string
	FullName string

	ErlName string
}

type MessageTypes []*MessageType

func (messageType *MessageType) FromDescriptor(fd *descriptor.FileDescriptorProto, d *descriptor.DescriptorProto, parent *MessageType) error {
	mt := MessageType{
		Parent: parent,

		Package: fd.GetPackage(),
		Name:    d.GetName(),
	}

	for p := parent; p != nil; p = p.Parent {
		mt.FullName += p.Name + "."
	}
	mt.FullName += mt.Name

	mt.ErlName = MessageTypeFullNameToErlRecordName(mt.FullName)

	*messageType = mt
	return nil
}

func MessageTypeFullNameToErlRecordName(name string) string {
	name2 := strings.ReplaceAll(name, ".", "__")
	return CamelCaseToSnakeCase(name2)
}
