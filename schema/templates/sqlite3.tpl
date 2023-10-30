{{ define "create_table" }}
CREATE TABLE {{- if .IfNotExists}} IF NOT EXISTS {{ end }} `{{.Name}}` (
  {{- range $i, $col := .Columns -}}
    {{- if $i}},{{end -}}
    {{- template "column" $col -}}
  {{- end}}
  
  {{- if gt .NumPrimary 1 -}},
PRIMARY KEY (
    {{- range $i, $col := .Columns -}}
      {{- if $col.IsPrimary -}}
        {{- if $i}}, {{end -}}
        '{{$col.Name}}'
      {{- end -}}
    {{end -}}
    )
  {{- end -}}

  {{- range $i, $col := .Columns -}}
    {{ if $col.ReferenceTo -}},
    FOREIGN KEY ('{{ $col.Name }}') REFERENCES `{{$col.ReferenceTo.Table}}`('{{ $col.ReferenceTo.Column }}')
    {{ end -}}
  {{- end -}}
)
{{ end }}

{{ define "column" }}
'{{.Name}}' {{GetType .Kind .KindLen}}
{{- if .SoloPrimary }} PRIMARY KEY{{- end -}}
{{- if .IsAutoincrement }} AUTOINCREMENT{{- end -}}
{{- if not .SoloPrimary }}{{ if not .IsNull }} NOT NULL{{ else }} NULL{{- end -}}{{- end -}}
{{- if .IsUnique }} UNIQUE{{- end -}}
{{- GetDefault .Kind .DefaultVal -}}
{{ end }}

{{ define "create_index" }}
CREATE{{ if .Unique }} UNIQUE{{- end }} INDEX{{ if .IfNotExists }} IF NOT EXISTS{{ end }}
{{- if .Name }} `{{.Name}}` {{ else }} 'unq_{{.Table.Name}}_{{ join .Columns "_"}}'{{ end }} ON `{{.Table.Name}}`(
  {{- range $i, $col := .Columns -}}
  {{- if $i }}, {{ end -}}
  '{{- $col -}}'
  {{- end -}}
)
{{ end }}
