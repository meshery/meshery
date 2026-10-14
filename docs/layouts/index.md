# {{ .Title | default site.Title }}
{{ with site.Params.metadata.description | default site.Params.description }}
> {{ . }}
{{- end }}

{{ if .RawContent }}{{ .RenderShortcodes }}{{ else }}{{ site.Params.description }}{{ end }}