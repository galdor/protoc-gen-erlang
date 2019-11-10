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
	"errors"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type FieldTypeId string

const (
	FieldTypeIdBool     FieldTypeId = "bool"
	FieldTypeIdFloat    FieldTypeId = "float"
	FieldTypeIdDouble   FieldTypeId = "double"
	FieldTypeIdInt32    FieldTypeId = "int32"
	FieldTypeIdInt64    FieldTypeId = "int64"
	FieldTypeIdUInt32   FieldTypeId = "uint32"
	FieldTypeIdUInt64   FieldTypeId = "uint64"
	FieldTypeIdSInt32   FieldTypeId = "sint32"
	FieldTypeIdSInt64   FieldTypeId = "sint64"
	FieldTypeIdFixed32  FieldTypeId = "fixed32"
	FieldTypeIdFixed64  FieldTypeId = "fixed64"
	FieldTypeIdSFixed32 FieldTypeId = "sfixed32"
	FieldTypeIdSFixed64 FieldTypeId = "sfixed64"
	FieldTypeIdString   FieldTypeId = "string"
	FieldTypeIdBytes    FieldTypeId = "bytes"
	FieldTypeIdGroup    FieldTypeId = "group" // deprecated
	FieldTypeIdEnum     FieldTypeId = "enum"
	FieldTypeIdMessage  FieldTypeId = "message"
)

func (tid *FieldTypeId) FromProto(v descriptor.FieldDescriptorProto_Type) error {
	switch v {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		*tid = FieldTypeIdDouble
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		*tid = FieldTypeIdFloat
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		*tid = FieldTypeIdInt64
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		*tid = FieldTypeIdUInt64
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		*tid = FieldTypeIdInt32
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		*tid = FieldTypeIdFixed64
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		*tid = FieldTypeIdFixed32
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		*tid = FieldTypeIdBool
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		*tid = FieldTypeIdString
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		*tid = FieldTypeIdGroup
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		*tid = FieldTypeIdMessage
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		*tid = FieldTypeIdBytes
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		*tid = FieldTypeIdUInt32
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		*tid = FieldTypeIdEnum
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		*tid = FieldTypeIdSFixed32
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		*tid = FieldTypeIdSFixed64
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		*tid = FieldTypeIdSInt32
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		*tid = FieldTypeIdSInt64
	default:
		return errors.New("invalid value")
	}

	return nil
}
