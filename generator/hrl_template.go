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
	"text/template"
)

var erlHRLTemplateContent = `
{{- define "erl_field" }}
  {{ .ErlName }} :: undefined | {{ .ErlTypeName }}
{{- end }}

{{- define "erl_message" }}
%% Generated for message type {{ .FullName }}.
-record({{ .ErlName }}, {
  {{- range $i, $f := .Fields }}
  {{- if gt $i 0 }},{{ end }}{{- template "erl_field" . }}
  {{- end }}
}).
{{- end }}

%%% Generated from protobuf package {{ .PackageName }}.
%%% DO NOT EDIT.

{{ range .PackageMessageTypes }}
{{- template "erl_message" . }}
{{ end }}
`

func ErlHRLTemplate() (*template.Template, error) {
	tpl := template.New("erl_hrl")

	if _, err := tpl.Parse(erlHRLTemplateContent); err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	return tpl, nil
}
