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
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

var verbose = true

func main() {
	// Read the request
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		die("cannot read stdin: %v", err)
	}

	var req plugin.CodeGeneratorRequest
	if err := proto.Unmarshal(input, &req); err != nil {
		die("cannot decode request: %v", err)
	}

	// Process input files
	var res plugin.CodeGeneratorResponse
	// TODO

	// Write the response
	output, err := proto.Marshal(&res)
	if err != nil {
		die("cannot encode response: %v", err)
	}

	if _, err := os.Stdout.Write(output); err != nil {
		die("cannot write to stdout: %v", err)
	}
}

func info(format string, args ...interface{}) {
	if !verbose {
		return
	}

	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
