package templates

const LearningPathTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "learning-path"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
id: "{{ .ID }}"
weight: {{ .Weight }}
---
`

const CourseTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "course"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`

const ModuleTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "module"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`

const PageTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "page"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`

const LabTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "lab"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`

const TestTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "test"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`

const CertificationTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "certification"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
id: "{{ .ID }}"
weight: {{ .Weight }}
---
`

const ExamTemplate = `---
title: "{{ .Title }}"
description: "{{ .Description }}"
type: "exam"
level: "{{ .Level }}"
{{- if .Tags }}
tags:
{{- range .Tags }}
  - {{ . }}
{{- end }}
{{- else }}
tags: []
{{- end }}
{{- if .Category }}
categories: "{{ .Category }}"
{{- end }}
weight: {{ .Weight }}
---
`
