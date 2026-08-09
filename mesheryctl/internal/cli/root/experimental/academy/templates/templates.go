package templates

const LearningPathTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "learning-path"
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
id: {{ .ID | yamlQuote }}
weight: {{ .Weight }}
---
`

const CourseTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "course"
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
weight: {{ .Weight }}
---
`

const ModuleTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "module"
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
weight: {{ .Weight }}
---
`

const PageTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "page"
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
weight: {{ .Weight }}
---
`

const LabTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "lab"
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
weight: {{ .Weight }}
---
`

const TestTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "test"
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
weight: {{ .Weight }}
---
`

const CertificationTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "certification"
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
id: {{ .ID | yamlQuote }}
weight: {{ .Weight }}
---
`

const ExamTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "exam"
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
weight: {{ .Weight }}
---
`

const ChallengeTemplate = `---
title: {{ .Title | yamlQuote }}
description: {{ .Description | yamlQuote }}
type: "challenge"
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
id: {{ .ID | yamlQuote }}
weight: {{ .Weight }}
---
`
