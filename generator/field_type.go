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
	OneofType   *OneofType   // for oneof fields

	TypeId   FieldTypeId
	TypeName string

	ErlName          string
	ErlValueTypeSpec string // available after type resolution
	ErlTypeSpec      string // available after type resolution
	ErlDefaultValue  string // available after type resolution
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

	if fid.GetDefaultValue() != "" {
		return fmt.Errorf("unsupported default value")
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
		ft.ErlValueTypeSpec = "boolean()"
		ft.ErlDefaultValue = "false"
	case FieldTypeIdFloat:
		ft.ErlValueTypeSpec = "float()"
		ft.ErlDefaultValue = "0.0"
	case FieldTypeIdDouble:
		ft.ErlValueTypeSpec = "float()"
		ft.ErlDefaultValue = "0.0"
	case FieldTypeIdInt32:
		ft.ErlValueTypeSpec = "-2147483648..2147483647"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdInt64:
		ft.ErlValueTypeSpec = "-9223372036854775808..9223372036854775807"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdUInt32:
		ft.ErlValueTypeSpec = "0..4294967295"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdUInt64:
		ft.ErlValueTypeSpec = "0..18446744073709551615"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdSInt32:
		ft.ErlValueTypeSpec = "-2147483648..2147483647"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdSInt64:
		ft.ErlValueTypeSpec = "-9223372036854775808..9223372036854775807"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdFixed32:
		ft.ErlValueTypeSpec = "-2147483648..2147483647"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdFixed64:
		ft.ErlValueTypeSpec = "-9223372036854775808..9223372036854775807"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdSFixed32:
		ft.ErlValueTypeSpec = "-2147483648..2147483647"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdSFixed64:
		ft.ErlValueTypeSpec = "-9223372036854775808..9223372036854775807"
		ft.ErlDefaultValue = "0"
	case FieldTypeIdString:
		ft.ErlValueTypeSpec = "iodata()"
		ft.ErlDefaultValue = "[]"
	case FieldTypeIdBytes:
		ft.ErlValueTypeSpec = "iodata()"
		ft.ErlDefaultValue = "[]"

	case FieldTypeIdGroup:
		return fmt.Errorf("unsupported field type %q", ft.TypeId)

	case FieldTypeIdEnum:
		et := absNameResolver.FindEnumType(ft.TypeName)
		if et == nil {
			return fmt.Errorf("unknown enum type %q", ft.TypeName)
		}

		ft.EnumType = et
		ft.ErlValueTypeSpec = et.ErlPackage + ":" + et.ErlName + "()"

		ft.ErlDefaultValue = et.Values[0].ErlName

	case FieldTypeIdMessage:
		mt := absNameResolver.FindMessageType(ft.TypeName)
		if mt == nil {
			return fmt.Errorf("unknown message type %q",
				ft.TypeName)
		}

		ft.MessageType = mt
		ft.ErlValueTypeSpec = mt.ErlPackage + ":" + mt.ErlName + "()"
		ft.ErlDefaultValue = "undefined"

	default:
		return fmt.Errorf("unhandled type %q", string(ft.TypeId))
	}

	ft.ErlTypeSpec = ft.ErlValueTypeSpec

	if ft.Repeated {
		ft.ErlTypeSpec = "list(" + ft.ErlTypeSpec + ")"
		ft.ErlDefaultValue = "[]"
	}

	if ft.TypeId == FieldTypeIdMessage && !ft.Repeated {
		ft.ErlTypeSpec = "undefined | " + ft.ErlTypeSpec
	}

	return nil
}
