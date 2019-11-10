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

var erlModuleTemplateContent = `
{{- define "erl_enum" }}
%% Generated for enum type {{ .FullName }}.
-type {{ .ErlName }}() ::{{ range $i, $v := .Values }}{{ if gt $i 0 }} |{{ end}} {{ .ErlName }}{{ end }}.
{{- end }}

{{- define "erl_message" }}
%% Generated for message type {{ .FullName }}.
-type {{ .ErlName }}() :: #{{ .ErlName }}{}.
{{- end }}

%%% Generated from protobuf package {{ .PackageName }}.
%%% DO NOT EDIT.

-module({{ .ErlModuleName }}).

-include("{{ .ErlModuleName }}.hrl").

-export_type([
  {{- range $i, $e := .PackageEnumTypes }}
  {{- if gt $i 0 }},{{ end }}
  {{ $e.ErlName }}/0
  {{- end }}
]).

-export_type([
  {{- range $i, $m := .PackageMessageTypes }}
  {{- if gt $i 0 }},{{ end }}
  {{ $m.ErlName }}/0
  {{- end }}
]).

{{ range .PackageEnumTypes }}
{{ template "erl_enum" . }}
{{ end }}

{{ range .PackageMessageTypes }}
{{ template "erl_message" . }}
{{ end }}
`

func ErlModuleTemplate() (*template.Template, error) {
	tpl := template.New("erl_module")

	if _, err := tpl.Parse(erlModuleTemplateContent); err != nil {
		return nil, fmt.Errorf("cannot parse template: %w", err)
	}

	return tpl, nil
}
