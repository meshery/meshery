apiVersion: meshery.layer5.io/v1alpha1
kind: Broker
metadata:
  name: {{ .Values.name }}
  namespace: {{ .Release.Namespace }}
spec:
  size: {{ .Values.replica }}
