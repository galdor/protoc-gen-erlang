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

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type FieldType struct {
	Name   string
	Number int

	Repeated bool
	Required bool
	Optional bool

	EnumType    *EnumType    // for enum fields
	MessageType *MessageType // for message fields

	TypeId   FieldTypeId
	TypeName string

	// TODO default value (representation?)

	ErlName     string
	ErlTypeName string // available after type resolution
}

type FieldTypes []*FieldType

func (fieldType *FieldType) FromDescriptor(fid *descriptor.FieldDescriptorProto) error {
	ft := FieldType{
		Name:   fid.GetName(),
		Number: int(fid.GetNumber()),

		TypeName: fid.GetTypeName(),
	}

	switch fid.GetLabel() {
	case descriptor.FieldDescriptorProto_LABEL_REPEATED:
		ft.Repeated = true
	case descriptor.FieldDescriptorProto_LABEL_REQUIRED:
		ft.Required = true
	case descriptor.FieldDescriptorProto_LABEL_OPTIONAL:
		ft.Optional = true
	default:
		return fmt.Errorf("unsupported field label %q", fid.Label)
	}

	if fid.GetExtendee() != "" {
		return fmt.Errorf("unsupported extendee field")
	}

	if err := ft.TypeId.FromProto(fid.GetType()); err != nil {
		return fmt.Errorf("invalid type %d: %w", fid.GetType(), err)
	}

	ft.ErlName = ft.Name

	*fieldType = ft
	return nil
}

func (ft *FieldType) ResolveType(absNameResolver AbsoluteNameResolver) error {
	switch ft.TypeId {
	case FieldTypeIdBool:
		ft.ErlTypeName = "boolean()"
	case FieldTypeIdFloat:
		ft.ErlTypeName = "float()"
	case FieldTypeIdDouble:
		ft.ErlTypeName = "float()"
	case FieldTypeIdInt32:
		ft.ErlTypeName = "-2147483648..2147483647"
	case FieldTypeIdInt64:
		ft.ErlTypeName = "-9223372036854775808..9223372036854775807"
	case FieldTypeIdUInt32:
		ft.ErlTypeName = "0..4294967295"
	case FieldTypeIdUInt64:
		ft.ErlTypeName = "0..18446744073709551615"
	case FieldTypeIdSInt32:
		ft.ErlTypeName = "-2147483648..2147483647"
	case FieldTypeIdSInt64:
		ft.ErlTypeName = "-9223372036854775808..9223372036854775807"
	case FieldTypeIdFixed32:
		ft.ErlTypeName = "-2147483648..2147483647"
	case FieldTypeIdFixed64:
		ft.ErlTypeName = "-9223372036854775808..9223372036854775807"
	case FieldTypeIdSFixed32:
		ft.ErlTypeName = "-2147483648..2147483647"
	case FieldTypeIdSFixed64:
		ft.ErlTypeName = "-9223372036854775808..9223372036854775807"
	case FieldTypeIdString:
		ft.ErlTypeName = "iodata()"
	case FieldTypeIdBytes:
		ft.ErlTypeName = "iodata()"

	case FieldTypeIdGroup:
		return fmt.Errorf("unsupported field type %q", ft.TypeId)

	case FieldTypeIdEnum:
		et := absNameResolver.FindEnumType(ft.TypeName)
		if et == nil {
			return fmt.Errorf("unknown enum type %q", ft.TypeName)
		}

		ft.EnumType = et
		ft.ErlTypeName = et.ErlPackage + ":" + et.ErlName + "()"

	case FieldTypeIdMessage:
		mt := absNameResolver.FindMessageType(ft.TypeName)
		if mt == nil {
			return fmt.Errorf("unknown message type %q",
				ft.TypeName)
		}

		ft.MessageType = mt
		ft.ErlTypeName = mt.ErlPackage + ":" + mt.ErlName + "()"

	default:
		return fmt.Errorf("unhandled type %q", string(ft.TypeId))
	}

	if ft.Repeated {
		ft.ErlTypeName = "list(" + ft.ErlTypeName + ")"
	}

	return nil
}
