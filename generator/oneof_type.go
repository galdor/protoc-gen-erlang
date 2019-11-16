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

type OneofType struct {
	Message *MessageType

	Name string

	Fields FieldTypes

	ErlName         string
	ErlTypeSpec     string // available after type resolution
	ErlDefaultValue string // available after type resolution
}

type OneofTypes []*OneofType

func (oneofType *OneofType) FromDescriptor(od *descriptor.OneofDescriptorProto, msg *MessageType) error {
	ot := OneofType{
		Message: msg,
		Name:    od.GetName(),
	}

	ot.ErlName = ot.Name

	*oneofType = ot
	return nil
}

func (ot *OneofType) AddField(ft *FieldType) {
	ot.Fields = append(ot.Fields, ft)
}

func (ot *OneofType) ResolveType(absNameResolver AbsoluteNameResolver) error {
	typeSpecs := make([]string, len(ot.Fields))

	for i, ft := range ot.Fields {
		typeSpecs[i] = fmt.Sprintf("{%s, %s}",
			ft.ErlName, ft.ErlValueTypeSpec)
	}

	ot.ErlTypeSpec = "undefined | " + strings.Join(typeSpecs, " | ")
	ot.ErlDefaultValue = "undefined"

	return nil
}

func OneofTypeNameToErlName(name string, msg *MessageType) string {
	return fmt.Sprintf("%s_%s")
}
