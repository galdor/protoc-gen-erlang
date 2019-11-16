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
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type MessageType struct {
	Parent *MessageType

	Package      string
	Name         string
	FullName     string
	AbsoluteName string

	ErlPackage string
	ErlName    string

	Oneofs OneofTypes
	Fields FieldTypes
}

type MessageTypes []*MessageType

func (messageType *MessageType) FromDescriptor(fd *descriptor.FileDescriptorProto, d *descriptor.DescriptorProto, parent *MessageType) error {
	mt := MessageType{
		Parent: parent,

		Package: fd.GetPackage(),
		Name:    d.GetName(),
	}

	mt.FullName = MessageTypeFullName(&mt)
	mt.AbsoluteName = "." + mt.Package + "." + mt.FullName

	mt.ErlPackage = ProtoPackageNameToErlModuleName(mt.Package)
	mt.ErlName = MessageTypeFullNameToErlRecordName(mt.FullName)

	for _, od := range d.OneofDecl {
		var ot OneofType
		if err := ot.FromDescriptor(od, &mt); err != nil {
			return fmt.Errorf("invalid oneof %q: %w",
				od.GetName(), err)
		}

		mt.Oneofs = append(mt.Oneofs, &ot)
	}

	for _, fid := range d.Field {
		var ft FieldType
		if err := ft.FromDescriptor(fid); err != nil {
			return fmt.Errorf("invalid field %q: %w",
				fid.GetName(), err)
		}

		if fid.OneofIndex != nil {
			idx := fid.GetOneofIndex()
			if int(idx) >= len(mt.Oneofs) {
				return fmt.Errorf("invalid index %d for "+
					"field %q", idx, ft.Name)
			}

			ot := mt.Oneofs[idx]
			ot.AddField(&ft)

			ft.OneofType = ot
		}

		mt.Fields = append(mt.Fields, &ft)
	}

	*messageType = mt
	return nil
}

func (mt *MessageType) ResolveTypes(absNameResolver AbsoluteNameResolver) error {
	for _, ft := range mt.Fields {
		if err := ft.ResolveType(absNameResolver); err != nil {
			return fmt.Errorf("cannot resolve type of field %q: %w",
				ft.Name, err)
		}
	}

	for _, ot := range mt.Oneofs {
		if err := ot.ResolveType(absNameResolver); err != nil {
			return fmt.Errorf("cannot resolve type of oneof %q: %w",
				ot.Name, err)
		}
	}

	return nil
}

func MessageTypeFullName(mt *MessageType) string {
	var parts []string

	for ; mt != nil; mt = mt.Parent {
		parts = append(parts, mt.Name)
	}

	for i := len(parts)/2 - 1; i >= 0; i-- {
		j := len(parts) - 1 - i
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}

func MessageTypeFullNameToErlRecordName(name string) string {
	name2 := strings.ReplaceAll(name, ".", "_")
	return CamelCaseToSnakeCase(name2)
}
