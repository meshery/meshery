# {{ .Title }}
{{ with .Params.description }}
> {{ . }}
{{- end }}

{{ if .RawContent }}{{ .RenderShortcodes }}{{ else }}{{ .Params.description }}{{ end }}