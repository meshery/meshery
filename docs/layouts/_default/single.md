# {{ .Title }}
{{ with .Date }}
date: {{ .Format "2006-01-02" }}
{{- end }}
{{ with .Params.description }}
> {{ . }}
{{- end }}

{{ if .RawContent }}{{ .RenderShortcodes }}{{ else }}{{ .Summary | plainify }}{{ end }}