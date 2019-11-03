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
)

func CamelCaseToSnakeCase(s string) string {
	var buf bytes.Buffer

	isUpper := func(c byte) bool { return c >= 'A' && c <= 'Z' }
	isLower := func(c byte) bool { return c >= 'a' && c <= 'z' }
	toLower := func(c byte) byte { return c + 'a' - 'A' }

	cs := []byte(s)

	for i, c := range cs {
		if isLower(c) {
			buf.WriteByte(c)
			continue
		}

		wordStart := i > 0 && isLower(cs[i-1])
		wordEnd := i < len(cs)-1 && isLower(cs[i+1])

		if i > 0 && (wordStart || wordEnd) && c != '_' && cs[i-1] != '_' {
			buf.WriteByte('_')
		}

		if isUpper(c) {
			c = toLower(c)
		}

		buf.WriteByte(c)
	}

	return buf.String()
}
