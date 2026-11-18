{{- $svc := .Values.service | default dict }}
{{- $type := $svc.type | default "" }}
{{- if and $type (not (has $type (list "ClusterIP" "NodePort" "LoadBalancer"))) }}
{{- fail (printf "meshery-broker.service.type %q is not one of ClusterIP, NodePort, LoadBalancer" $type) }}
{{- end }}
{{- /* Stricter than the CRD's CEL rules, which admit these fields when type is
       unset; the operator then silently ignores them on a ClusterIP Service. */}}
{{- if and (ne $type "LoadBalancer") (or $svc.loadBalancerClass $svc.loadBalancerSourceRanges) }}
{{- fail "meshery-broker.service.loadBalancerClass and loadBalancerSourceRanges require meshery-broker.service.type: LoadBalancer" }}
{{- end }}
{{- /* v1alpha1, not the v1alpha2 storage version: the schemas are
       field-identical (conversion None), Meshery Server applies v1alpha1, and
       Helm resolves a manifest's apiVersion before the pre-upgrade CRD Job
       runs, so v1alpha2 would fail upgrades on clusters whose CRDs predate it. */}}
apiVersion: meshery.io/v1alpha1
kind: Broker
metadata:
  name: {{ .Values.name }}
  namespace: {{ .Release.Namespace }}
spec:
  size: {{ .Values.replica }}
  {{- with .Values.version }}
  version: {{ . | quote }}
  {{- end }}
  {{- if or $type $svc.annotations $svc.loadBalancerClass $svc.loadBalancerSourceRanges $svc.externalEndpointOverride }}
  service:
    {{- with $type }}
    type: {{ . }}
    {{- end }}
    {{- with $svc.annotations }}
    annotations:
      {{- range $k, $v := . }}
      {{ $k | quote }}: {{ $v | quote }}
      {{- end }}
    {{- end }}
    {{- with $svc.loadBalancerClass }}
    loadBalancerClass: {{ . | quote }}
    {{- end }}
    {{- with $svc.loadBalancerSourceRanges }}
    loadBalancerSourceRanges:
      {{- range . }}
      - {{ . | quote }}
      {{- end }}
    {{- end }}
    {{- with $svc.externalEndpointOverride }}
    externalEndpointOverride: {{ . | quote }}
    {{- end }}
  {{- end }}
