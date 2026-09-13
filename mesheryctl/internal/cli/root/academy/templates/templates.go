package templates

const NodeTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
{{- if .BannerString }}
banner: {{ .BannerString | yamlQuote }}
{{- end }}
type: {{ .TypeString | yamlQuote }}
level: {{ .LevelString | yamlQuote }}
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
weight: {{ .WeightInt }}
{{- if .Draft }}
draft: true
{{- end }}
---
`
