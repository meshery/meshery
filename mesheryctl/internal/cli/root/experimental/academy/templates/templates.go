package templates

const NodeTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: {{ .Type | yamlQuote }}
level: {{ .Level | yamlQuote }}
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . | yamlQuote }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: {{ .Category | yamlQuote }}
{{- end }}
{{- if .ID }}
id: {{ .ID | yamlQuote }}
{{- end }}
weight: {{ .Weight }}
---
`
