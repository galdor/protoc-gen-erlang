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
	"os"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type Generator struct {
	Request  *plugin.CodeGeneratorRequest
	Response *plugin.CodeGeneratorResponse

	Verbose bool
}

func NewGenerator(req *plugin.CodeGeneratorRequest) (*Generator, error) {
	g := Generator{
		Request:  req,
		Response: new(plugin.CodeGeneratorResponse),

		Verbose: true,
	}

	return &g, nil
}

func (g *Generator) Info(format string, args ...interface{}) {
	if !g.Verbose {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func (g *Generator) GenerateOutput() error {
	// TODO
	return nil
}
