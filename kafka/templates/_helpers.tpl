{{/*
Copyright Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/* vim: set filetype=mustache: */}}

{{/*
Return the proper Kafka controller-eligible fullname
*/}}
{{- define "kafka.controller.fullname" -}}
{{- printf "%s-controller" (include "common.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the proper Kafka broker fullname
*/}}
{{- define "kafka.broker.fullname" -}}
{{- printf "%s-broker" (include "common.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
 Create the name of the service account to use
 */}}
{{- define "kafka.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "common.names.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Kafka image name
*/}}
{{- define "kafka.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name (for the init container volume-permissions image)
*/}}
{{- define "kafka.volumePermissions.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.defaultInitContainers.volumePermissions.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name (for the init container auto-discovery image)
*/}}
{{- define "kafka.autoDiscovery.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.defaultInitContainers.autoDiscovery.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper JMX exporter image name
*/}}
{{- define "kafka.metrics.jmx.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.metrics.jmx.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "kafka.imagePullSecrets" -}}
{{ include "common.images.pullSecrets" (dict "images" (list .Values.image .Values.defaultInitContainers.volumePermissions.image .Values.defaultInitContainers.autoDiscovery.image .Values.metrics.jmx.image) "global" .Values.global) }}
{{- end -}}

{{/*
Return true if encryption via TLS for client connections should be configured
*/}}
{{- define "kafka.sslEnabled" -}}
{{- $res := "" -}}
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker .Values.listeners.controller -}}
{{- range $i := .Values.listeners.extraListeners -}}
{{- $listeners = append $listeners $i -}}
{{- end -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SSL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if SASL connections should be configured
*/}}
{{- define "kafka.saslEnabled" -}}
{{- $res := "" -}}
{{- if include "kafka.client.saslEnabled" . -}}
{{- $res = "true" -}}
{{- else -}}
{{- $listeners := list .Values.listeners.interbroker .Values.listeners.controller -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SASL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if SASL connections should be configured
*/}}
{{- define "kafka.client.saslEnabled" -}}
{{- $res := "" -}}
{{- $listeners := list .Values.listeners.client -}}
{{- range $i := .Values.listeners.extraListeners -}}
{{- $listeners = append $listeners $i -}}
{{- end -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SASL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Returns true if a SASL mechanism that uses usernames and passwords is in use
*/}}
{{- define "kafka.saslUserPasswordsEnabled" -}}
{{- if (include "kafka.saslEnabled" .) -}}
{{- if or (regexFind "PLAIN" (upper .Values.sasl.enabledMechanisms)) (regexFind "SCRAM" (upper .Values.sasl.enabledMechanisms)) -}}
true
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Returns true if a SASL mechanism that uses client IDs and client secrets is in use
*/}}
{{- define "kafka.saslClientSecretsEnabled" -}}
{{- if (include "kafka.saslEnabled" .) -}}
{{- if (regexFind "OAUTHBEARER" (upper .Values.sasl.enabledMechanisms)) -}}
true
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Returns the security module based on the provided sasl mechanism
*/}}
{{- define "kafka.saslSecurityModule" -}}
{{- if eq "PLAIN" .mechanism -}}
org.apache.kafka.common.security.plain.PlainLoginModule required
{{- else if regexFind "SCRAM" .mechanism -}}
org.apache.kafka.common.security.scram.ScramLoginModule required
{{- else if eq "OAUTHBEARER" .mechanism -}}
org.apache.kafka.common.security.oauthbearer.OAuthBearerLoginModule required
{{- end -}}
{{- end -}}

{{/*
Return the Kafka Kraft secret
*/}}
{{- define "kafka.kraftSecretName" -}}
{{- if .Values.existingKraftSecret -}}
    {{- print (tpl .Values.existingKraftSecret .) -}}
{{- else -}}
    {{- printf "%s-kraft" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka SASL credentials secret
*/}}
{{- define "kafka.saslSecretName" -}}
{{- if .Values.sasl.existingSecret -}}
    {{- print (tpl .Values.sasl.existingSecret .) -}}
{{- else -}}
    {{- printf "%s-user-passwords" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a SASL credentials secret object should be created
*/}}
{{- define "kafka.createSaslSecret" -}}
{{- if and (include "kafka.saslEnabled" .) (empty .Values.sasl.existingSecret) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.tlsSecretName" -}}
{{- if .Values.tls.existingSecret -}}
    {{- print (tpl .Values.tls.existingSecret .) -}}
{{- else -}}
    {{- printf "%s-tls" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.createTlsSecret" -}}
{{- if and (include "kafka.sslEnabled" .) (empty .Values.tls.existingSecret) .Values.tls.autoGenerated.enabled (eq .Values.tls.autoGenerated.engine "helm") -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a Certificate object should be created
*/}}
{{- define "kafka.createCertificate" -}}
{{- if and (include "kafka.sslEnabled" .) (empty .Values.tls.existingSecret) .Values.tls.autoGenerated.enabled (eq .Values.tls.autoGenerated.engine "cert-manager") -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka TLS credentials secret
*/}}
{{- define "kafka.tlsPasswordsSecretName" -}}
{{- if .Values.tls.passwordsSecret -}}
    {{- print (tpl .Values.tls.passwordsSecret .) -}}
{{- else -}}
    {{- printf "%s-tls-passwords" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.createTlsPasswordsSecret" -}}
{{- if and (include "kafka.sslEnabled" .) (or (empty .Values.tls.passwordsSecret) .Values.tls.autoGenerated.enabled) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Returns the secret name for the Kafka Provisioning client
*/}}
{{- define "kafka.client.passwordsSecretName" -}}
{{- if .Values.provisioning.auth.tls.passwordsSecret -}}
    {{- print (tpl .Values.provisioning.auth.tls.passwordsSecret .) -}}
{{- else -}}
    {{- printf "%s-client-secret" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use for the Kafka Provisioning client
*/}}
{{- define "kafka.provisioning.serviceAccountName" -}}
{{- if .Values.provisioning.serviceAccount.create -}}
    {{ default (printf "%s-provisioning" (include "common.names.fullname" .)) .Values.provisioning.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.provisioning.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible configuration configmap
*/}}
{{- define "kafka.controller.configmapName" -}}
{{- if .Values.controller.existingConfigmap -}}
    {{- print (tpl .Values.controller.existingConfigmap .) -}}
{{- else if .Values.existingConfigmap -}}
    {{- print (tpl .Values.existingConfigmap .) -}}
{{- else -}}
    {{- printf "%s-configuration" (include "kafka.controller.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible secret configuration
*/}}
{{- define "kafka.controller.secretConfigName" -}}
{{- if .Values.controller.existingSecretConfig -}}
    {{- print (tpl .Values.controller.existingSecretConfig .) -}}
{{- else if .Values.existingSecretConfig -}}
    {{- print (tpl .Values.controller.existingSecretConfig .) -}}
{{- else -}}
    {{- printf "%s-secret-configuration" (include "kafka.controller.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible secret configuration values
*/}}
{{- define "kafka.controller.secretConfig" -}}
{{- if .Values.secretConfig }}
    {{- print (tpl .Values.secretConfig .) -}}
{{- end }}
{{- if .Values.controller.secretConfig }}
    {{- print (tpl .Values.controller.secretConfig .) -}}
{{- end }}
{{- end -}}

{{/*
Return true if a configmap object should be created for controller-eligible pods
*/}}
{{- define "kafka.controller.createConfigmap" -}}
{{- if and (not .Values.controller.existingConfigmap) (not .Values.existingConfigmap) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config should be created for controller-eligible pods
*/}}
{{- define "kafka.controller.createSecretConfig" -}}
{{- if and (or .Values.controller.secretConfig .Values.secretConfig) (and (not .Values.controller.existingSecretConfig) (not .Values.existingSecretConfig)) }}
    {{- true -}}
{{- end -}}
{{- end -}}
{{/*
Return true if a secret object with config exists for controller-eligible pods
*/}}
{{- define "kafka.controller.secretConfigExists" -}}
{{- if or .Values.controller.secretConfig .Values.secretConfig .Values.controller.existingSecretConfig .Values.existingSecretConfig }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker configuration configmap
*/}}
{{- define "kafka.broker.configmapName" -}}
{{- if .Values.broker.existingConfigmap -}}
    {{- print (tpl .Values.broker.existingConfigmap .) -}}
{{- else if .Values.existingConfigmap -}}
    {{- print (tpl .Values.existingConfigmap .) -}}
{{- else -}}
    {{- printf "%s-configuration" (include "kafka.broker.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker secret configuration
*/}}
{{- define "kafka.broker.secretConfigName" -}}
{{- if .Values.broker.existingSecretConfig -}}
    {{- print (tpl .Values.broker.existingSecretConfig .) -}}
{{- else if .Values.existingSecretConfig -}}
    {{- print (tpl .Values.existingSecretConfig .) -}}
{{- else -}}
    {{- printf "%s-secret-configuration" (include "kafka.broker.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker secret configuration values
*/}}
{{- define "kafka.broker.secretConfig" -}}
{{- if .Values.secretConfig }}
    {{- print (tpl .Values.secretConfig .) -}}
{{- end }}
{{- if .Values.broker.secretConfig }}
    {{- print (tpl .Values.broker.secretConfig .) -}}
{{- end }}
{{- end -}}

{{/*
Return true if a configmap object should be created for broker pods
*/}}
{{- define "kafka.broker.createConfigmap" -}}
{{- if and (not .Values.broker.existingConfigmap) (not .Values.existingConfigmap) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config should be created for broker pods
*/}}
{{- define "kafka.broker.createSecretConfig" -}}
{{- if and (or .Values.broker.secretConfig .Values.secretConfig) (and (not .Values.broker.existingSecretConfig) (not .Values.existingSecretConfig)) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config exists for broker pods
*/}}
{{- define "kafka.broker.secretConfigExists" -}}
{{- if or .Values.broker.secretConfig .Values.secretConfig .Values.broker.existingSecretConfig .Values.existingSecretConfig }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka log4j2 ConfigMap name.
*/}}
{{- define "kafka.log4j2.configMapName" -}}
{{- if .Values.existingLog4j2ConfigMap -}}
    {{- print (tpl .Values.existingLog4j2ConfigMap .) -}}
{{- else -}}
    {{- printf "%s-log4j2-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka configuration configmap
*/}}
{{- define "kafka.metrics.jmx.configmapName" -}}
{{- if .Values.metrics.jmx.existingConfigmap -}}
    {{- print (tpl .Values.metrics.jmx.existingConfigmap .) -}}
{{- else -}}
    {{ printf "%s-jmx-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a configmap object should be created
*/}}
{{- define "kafka.metrics.jmx.createConfigmap" -}}
{{- if and .Values.metrics.jmx.enabled .Values.metrics.jmx.config (not .Values.metrics.jmx.existingConfigmap) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Returns the Kafka listeners settings based on the listeners.* object
*/}}
{{- define "kafka.listeners" -}}
{{- if .context.Values.listeners.overrideListeners -}}
  {{- print .context.Values.listeners.overrideListeners -}}
{{- else -}}
  {{- $listeners := list .context.Values.listeners.client .context.Values.listeners.interbroker -}}
  {{- if .context.Values.externalAccess.enabled -}}
  {{- $listeners = append $listeners .context.Values.listeners.external -}}
  {{- end -}}
  {{- if .isController -}}
    {{- if .context.Values.controller.controllerOnly -}}
      {{- $listeners = list .context.Values.listeners.controller -}}
    {{- else -}}
      {{- $listeners = append $listeners .context.Values.listeners.controller -}}
      {{- range $i := .context.Values.listeners.extraListeners -}}
      {{- $listeners = append $listeners $i -}}
      {{- end -}}
    {{- end -}}
  {{- else -}}
    {{- range $i := .context.Values.listeners.extraListeners -}}
    {{- $listeners = append $listeners $i -}}
    {{- end -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s://:%d" (upper $listener.name) (int $listener.containerPort)) -}}
  {{- end -}}
  {{- join "," $res -}}
{{- end -}}
{{- end -}}

{{/*
Returns the list of advertised listeners, although the advertised address will be replaced during each node init time
*/}}
{{- define "kafka.advertisedListeners" -}}
{{- if .Values.listeners.advertisedListeners -}}
  {{- print .Values.listeners.advertisedListeners -}}
{{- else -}}
  {{- $listeners := list .Values.listeners.client .Values.listeners.interbroker -}}
  {{- range $i := .Values.listeners.extraListeners -}}
  {{- $listeners = append $listeners $i -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s://advertised-address-placeholder:%d" (upper $listener.name) (int $listener.containerPort)) -}}
  {{- end -}}
  {{- join "," $res -}}
{{- end -}}
{{- end -}}

{{/*
Returns the value listener.security.protocol.map based on the values of 'listeners.*.protocol'
*/}}
{{- define "kafka.securityProtocolMap" -}}
{{- if .context.Values.listeners.securityProtocolMap -}}
  {{- print .context.Values.listeners.securityProtocolMap -}}
{{- else -}}
  {{- $listeners := list .context.Values.listeners.controller .context.Values.listeners.client .context.Values.listeners.interbroker -}}
  {{- if and .isController .context.Values.controller.controllerOnly -}}
    {{- $listeners = list .context.Values.listeners.controller  -}}
  {{- else -}}
    {{- range $i := .context.Values.listeners.extraListeners -}}
    {{- $listeners = append $listeners $i -}}
    {{- end -}}
  {{- end -}}
  {{- if .context.Values.externalAccess.enabled -}}
    {{- $listeners = append $listeners .context.Values.listeners.external -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s:%s" (upper $listener.name) (upper $listener.protocol)) -}}
  {{- end -}}
  {{ join "," $res }}
{{- end -}}
{{- end -}}

{{/*
Returns the containerPorts for listeners.extraListeners
*/}}
{{- define "kafka.extraListeners.containerPorts" -}}
{{- range $listener := .Values.listeners.extraListeners -}}
- name: {{ lower $listener.name}}
  containerPort: {{ $listener.containerPort }}
{{ end }}
{{- end -}}

{{/*
Returns the controller quorum bootstrap servers based on the number of controller-eligible nodes
*/}}
{{- define "kafka.controller.quorumBootstrapServers" -}}
{{- if .Values.controller.quorumBootstrapServers -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.controller.quorumBootstrapServers "context" $) -}}
{{- else -}}
  {{- $fullname := include "kafka.controller.fullname" . }}
  {{- $serviceName := printf "%s-headless" (include "kafka.controller.fullname" .) | trunc 63 | trimSuffix "-" }}
  {{- $releaseNamespace := include "common.names.namespace" . -}}
  {{- $clusterDomain := .Values.clusterDomain }}
  {{- $port := int .Values.listeners.controller.containerPort }}
  {{- $bootstrapServers := list -}}
  {{- if and (.Values.controller.autoscaling) (.Values.controller.autoscaling.hpa) (.Values.controller.autoscaling.hpa.enabled) -}}
    {{- range $i := until (int .Values.controller.autoscaling.hpa.maxReplicas) -}}
      {{- $nodeAddress := printf "%s-%d.%s.%s.svc.%s:%d" $fullname (int $i) $serviceName $releaseNamespace $clusterDomain $port -}}
      {{- $bootstrapServers = append $bootstrapServers $nodeAddress -}}
    {{- end -}}
  {{- else -}}
    {{- range $i := until (int .Values.controller.replicaCount) -}}
      {{- $nodeAddress := printf "%s-%d.%s.%s.svc.%s:%d" $fullname (int $i) $serviceName $releaseNamespace $clusterDomain $port -}}
      {{- if eq (int $.Values.kraftVersion) 0 }}
        {{- $nodeId := add (int $i) (int $.Values.controller.minId) -}}
        {{- $bootstrapServers = append $bootstrapServers (printf "%d@%s" $nodeId $nodeAddress ) -}}
      {{- else }}
        {{- $bootstrapServers = append $bootstrapServers $nodeAddress -}}
      {{- end }}
    {{- end -}}
  {{- end -}}
  {{- join "," $bootstrapServers -}}
{{- end -}}
{{- end -}}

{{/*
Section of the server.properties shared by both controller-eligible and broker nodes
*/}}
{{- define "kafka.commonConfig" -}}
controller.listener.names: {{ .Values.listeners.controller.name }}
{{- if eq (int .Values.kraftVersion) 0 }}
controller.quorum.voters: {{ include "kafka.controller.quorumBootstrapServers" . }}
{{- else }}
controller.quorum.bootstrap.servers: {{ include "kafka.controller.quorumBootstrapServers" . }}
{{- end }}
{{- if include "kafka.sslEnabled" . }}
# TLS configuration
ssl.keystore.type: JKS
ssl.truststore.type: JKS
ssl.keystore.location: /opt/bitnami/kafka/config/certs/kafka.keystore.jks
ssl.truststore.location: /opt/bitnami/kafka/config/certs/kafka.truststore.jks
ssl.client.auth: {{ .Values.tls.sslClientAuth }}
ssl.endpoint.identification.algorithm: {{ .Values.tls.endpointIdentificationAlgorithm }}
{{- end }}
{{- if (include "kafka.saslEnabled" .) }}
# Listeners SASL JAAS configuration
sasl.enabled.mechanisms: {{ upper .Values.sasl.enabledMechanisms }}
{{- if regexFind "SASL" (upper .Values.listeners.interbroker.protocol) }}
sasl.mechanism.inter.broker.protocol: {{ upper .Values.sasl.interBrokerMechanism }}
{{- end }}
{{- if regexFind "SASL" (upper .Values.listeners.controller.protocol) }}
sasl.mechanism.controller.protocol: {{ upper .Values.sasl.controllerMechanism }}
{{- end }}
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker .Values.listeners.controller }}
{{- range $i := .Values.listeners.extraListeners }}
{{- $listeners = append $listeners $i }}
{{- end }}
{{- if .Values.externalAccess.enabled }}
{{- $listeners = append $listeners .Values.listeners.external }}
{{- end }}
{{- range $listener := $listeners }}
  {{- if and $listener.sslClientAuth (regexFind "SSL" (upper $listener.protocol)) }}
listener.name.{{lower $listener.name}}.ssl.client.auth: {{ $listener.sslClientAuth }}
  {{- end }}
  {{- if regexFind "SASL" (upper $listener.protocol) }}
    {{- range $mechanism := splitList "," $.Values.sasl.enabledMechanisms }}
      {{- $securityModule := include "kafka.saslSecurityModule" (dict "mechanism" (upper $mechanism)) }}
      {{- if and (eq (upper $mechanism) "OAUTHBEARER") (or (eq $listener.name $.Values.listeners.interbroker.name) (eq $listener.name $.Values.listeners.controller.name)) }}
listener.name.{{lower $listener.name}}.oauthbearer.sasl.login.callback.handler.class: org.apache.kafka.common.security.oauthbearer.secured.OAuthBearerLoginCallbackHandler
      {{- end }}
      {{- $saslJaasConfig := list $securityModule }}
      {{- if eq $listener.name $.Values.listeners.interbroker.name }}
        {{- if (eq (upper $mechanism) "OAUTHBEARER") }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "clientId=\"%s\"" $.Values.sasl.interbroker.clientId) }}
          {{- $saslJaasConfig = append $saslJaasConfig (print "clientSecret=\"interbroker-client-secret-placeholder\"") }}
        {{- else }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "username=\"%s\"" $.Values.sasl.interbroker.user) }}
          {{- $saslJaasConfig = append $saslJaasConfig (print "password=\"interbroker-password-placeholder\"") }}
        {{- end }}
      {{- else if eq $listener.name $.Values.listeners.controller.name }}
        {{- if (eq (upper $mechanism) "OAUTHBEARER") }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "clientId=\"%s\"" $.Values.sasl.controller.clientId) }}
          {{- $saslJaasConfig = append $saslJaasConfig (print "clientSecret=\"controller-client-secret-placeholder\"") }}
        {{- else }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "username=\"%s\"" $.Values.sasl.controller.user) }}
          {{- $saslJaasConfig = append $saslJaasConfig (print "password=\"controller-password-placeholder\"") }}
        {{- end }}
      {{- end }}
      {{- if eq (upper $mechanism) "PLAIN" }}
        {{- if eq $listener.name $.Values.listeners.interbroker.name }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"interbroker-password-placeholder\"" $.Values.sasl.interbroker.user) }}
        {{- else if eq $listener.name $.Values.listeners.controller.name }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"controller-password-placeholder\"" $.Values.sasl.controller.user) }}
        {{- end }}
        {{- range $i, $user := $.Values.sasl.client.users }}
          {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"password-placeholder-%d\"" $user (int $i)) }}
        {{- end }}
      {{- end }}
listener.name.{{lower $listener.name}}.{{lower $mechanism}}.sasl.jaas.config: {{ printf "%s;" (join " " $saslJaasConfig) }}
      {{- if eq (upper $mechanism) "OAUTHBEARER" }}
listener.name.{{lower $listener.name}}.oauthbearer.sasl.server.callback.handler.class: org.apache.kafka.common.security.oauthbearer.secured.OAuthBearerValidatorCallbackHandler
      {{- end }}
    {{- end }}
  {{- end }}
{{- end }}
{{- if regexFind "OAUTHBEARER" .Values.sasl.enabledMechanisms }}
sasl.oauthbearer.token.endpoint.url: {{ .Values.sasl.oauthbearer.tokenEndpointUrl }}
sasl.oauthbearer.jwks.endpoint.url: {{ .Values.sasl.oauthbearer.jwksEndpointUrl }}
sasl.oauthbearer.expected.audience: {{ .Values.sasl.oauthbearer.expectedAudience }}
sasl.oauthbearer.sub.claim.name: {{ .Values.sasl.oauthbearer.subClaimName }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Environment variables required to configure SASL
*/}}
{{- define "kafka.saslEnv" -}}
{{- if and (include "kafka.client.saslEnabled" . ) (include "kafka.saslUserPasswordsEnabled" .) .Values.sasl.client.users }}
- name: KAFKA_CLIENT_USERS
  value: {{ join "," .Values.sasl.client.users | quote }}
{{- if .Values.usePasswordFiles }}
- name: KAFKA_CLIENT_PASSWORDS_FILE
  value: /opt/bitnami/kafka/config/secrets/client-passwords
{{- else }}
- name: KAFKA_CLIENT_PASSWORDS
  valueFrom:
    secretKeyRef:
      name: {{ include "kafka.saslSecretName" . }}
      key: client-passwords
{{- end }}
{{- end }}
{{- if regexFind "SASL" (upper .Values.listeners.interbroker.protocol) }}
{{- if include "kafka.saslUserPasswordsEnabled" . }}
- name: KAFKA_INTER_BROKER_USER
  value: {{ .Values.sasl.interbroker.user | quote }}
{{- if .Values.usePasswordFiles }}
- name: KAFKA_INTER_BROKER_PASSWORD_FILE
  value: /opt/bitnami/kafka/config/secrets/inter-broker-password
{{- else }}
- name: KAFKA_INTER_BROKER_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "kafka.saslSecretName" . }}
      key: inter-broker-password
{{- end }}
{{- end }}
{{- if include "kafka.saslClientSecretsEnabled" . }}
- name: KAFKA_INTER_BROKER_CLIENT_ID
  value: {{ .Values.sasl.interbroker.clientId | quote }}
{{- if .Values.usePasswordFiles }}
- name: KAFKA_INTER_BROKER_CLIENT_SECRET_FILE
  value: /opt/bitnami/kafka/config/secrets/inter-broker-client-secret
{{- else }}
- name: KAFKA_INTER_BROKER_CLIENT_SECRET
  valueFrom:
    secretKeyRef:
      name: {{ include "kafka.saslSecretName" . }}
      key: inter-broker-client-secret
{{- end }}
{{- end }}
{{- end }}
{{- if regexFind "SASL" (upper .Values.listeners.controller.protocol) }}
{{- if include "kafka.saslUserPasswordsEnabled" . }}
- name: KAFKA_CONTROLLER_USER
  value: {{ .Values.sasl.controller.user | quote }}
{{- if .Values.usePasswordFiles }}
- name: KAFKA_CONTROLLER_PASSWORD_FILE
  value: /opt/bitnami/kafka/config/secrets/controller-password
{{- else }}
- name: KAFKA_CONTROLLER_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "kafka.saslSecretName" . }}
      key: controller-password
{{- end }}
{{- end }}
{{- if include "kafka.saslClientSecretsEnabled" . }}
- name: KAFKA_CONTROLLER_CLIENT_ID
  value: {{ .Values.sasl.controller.clientId | quote }}
{{- if .Values.usePasswordFiles }}
- name: KAFKA_CONTROLLER_CLIENT_SECRET_FILE
  value: /opt/bitnami/kafka/config/secrets/controller-client-secret
{{- else }}
- name: KAFKA_CONTROLLER_CLIENT_SECRET
  valueFrom:
    secretKeyRef:
      name: {{ include "kafka.saslSecretName" . }}
      key: controller-client-secret
{{- end }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Environment variables shared by both controller-eligible and broker nodes
*/}}
{{- define "kafka.commonEnv" -}}
- name: BITNAMI_DEBUG
  value: {{ ternary "true" "false" (or .Values.image.debug .Values.diagnosticMode.enabled) | quote }}
- name: KAFKA_KRAFT_CLUSTER_ID
  valueFrom:
    secretKeyRef:
      name: {{ template "kafka.kraftSecretName" . }}
      key: cluster-id
{{- if and (include "kafka.saslEnabled" .) (or (regexFind "SCRAM" (upper .Values.sasl.enabledMechanisms)) (regexFind "SCRAM" (upper .Values.sasl.controllerMechanism)) (regexFind "SCRAM" (upper .Values.sasl.interBrokerMechanism))) }}
- name: KAFKA_KRAFT_BOOTSTRAP_SCRAM_USERS
  value: "true"
{{ include "kafka.saslEnv" . }}
{{- end }}
{{- if .Values.provisioning.enabled }}
{{- end }}
{{- if .Values.metrics.jmx.enabled }}
- name: JMX_PORT
  value: {{ .Values.metrics.jmx.kafkaJmxPort | quote }}
{{- end }}
{{- end -}}

{{/*
Check if there are rolling tags in the images
*/}}
{{- define "kafka.checkRollingTags" -}}
{{- include "common.warnings.rollingTag" .Values.image }}
{{- include "common.warnings.rollingTag" .Values.defaultInitContainers.volumePermissions.image }}
{{- include "common.warnings.rollingTag" .Values.defaultInitContainers.autoDiscovery.image }}
{{- include "common.warnings.rollingTag" .Values.metrics.jmx.image }}
{{- end -}}

{{/*
Compile all warnings into a single message, and call fail.
*/}}
{{- define "kafka.validateValues" -}}
{{- $messages := list -}}
{{- $messages := append $messages (include "kafka.validateValues.listener.protocols" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.controller.nodePortListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.broker.nodePortListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.controller.externalIPListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.broker.externalIPListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.brokerRackAwareness" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.domainSpecified" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceType" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessAutoDiscoveryRBAC" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessAutoDiscoveryIPsOrNames" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerIPs" "context" .)) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerNames" "context" .)) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerAnnotations" "context" . )) -}}
{{- $messages := append $messages (include "kafka.validateValues.saslMechanisms" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.tlsSecret" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.provisioning.tlsPasswords" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.missingController" .) -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message | fail -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - Authentication protocols for Kafka */}}
{{- define "kafka.validateValues.listener.protocols" -}}
{{- $authProtocols := list "PLAINTEXT" "SASL_PLAINTEXT" "SASL_SSL" "SSL" -}}
{{- if not .Values.listeners.securityProtocolMap -}}
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker .Values.listeners.controller -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- $error := false -}}
{{- range $listener := $listeners -}}
{{- if not (has (upper $listener.protocol) $authProtocols) -}}
{{- $error := true -}}
{{- end -}}
{{- end -}}
{{- if $error -}}
kafka: listeners.*.protocol
    Available authentication protocols are "PLAINTEXT" "SASL_PLAINTEXT" "SSL" "SASL_SSL"
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of controller-eligible replicas must be the same as NodePort list in controller-eligible external service */}}
{{- define "kafka.validateValues.controller.nodePortListLength" -}}
{{- $replicaCount := int .Values.controller.replicaCount -}}
{{- $nodePortListLength := len .Values.externalAccess.controller.service.nodePorts -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.controller.service.nodePorts -}}
{{- $nodePortListLengthEqualsReplicaCount := eq $nodePortListLength $replicaCount -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.controller.service.externalIPs -}}
{{- if and .Values.externalAccess.enabled (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq .Values.externalAccess.controller.service.type "NodePort") (or (and (not $nodePortListIsEmpty) (not $nodePortListLengthEqualsReplicaCount)) (and $nodePortListIsEmpty $externalIPListIsEmpty)) -}}
kafka: .Values.externalAccess.controller.service.nodePorts
    Number of controller-eligible replicas and externalAccess.controller.service.nodePorts array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length nodePorts = {{ $nodePortListLength }} - {{ $externalIPListIsEmpty }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of broker replicas must be the same as NodePort list in broker external service */}}
{{- define "kafka.validateValues.broker.nodePortListLength" -}}
{{- $replicaCount := int .Values.broker.replicaCount -}}
{{- $nodePortListLength := len .Values.externalAccess.broker.service.nodePorts -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.broker.service.nodePorts -}}
{{- $nodePortListLengthEqualsReplicaCount := eq $nodePortListLength $replicaCount -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.broker.service.externalIPs -}}
{{- if and .Values.externalAccess.enabled (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq .Values.externalAccess.broker.service.type "NodePort") (or (and (not $nodePortListIsEmpty) (not $nodePortListLengthEqualsReplicaCount)) (and $nodePortListIsEmpty $externalIPListIsEmpty)) -}}
kafka: .Values.externalAccess.broker.service.nodePorts
    Number of broker replicas and externalAccess.broker.service.nodePorts array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length nodePorts = {{ $nodePortListLength }} - {{ $externalIPListIsEmpty }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as externalIPs list */}}
{{- define "kafka.validateValues.controller.externalIPListLength" -}}
{{- $replicaCount := int .Values.controller.replicaCount -}}
{{- $externalIPListLength := len .Values.externalAccess.controller.service.externalIPs -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.controller.service.externalIPs -}}
{{- $externalIPListEqualsReplicaCount := eq $externalIPListLength $replicaCount -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.controller.service.nodePorts -}}
{{- if and .Values.externalAccess.enabled (or .Values.externalAccess.controller.forceExpose (not .Values.controller.controllerOnly)) (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq .Values.externalAccess.controller.service.type "NodePort") (or (and (not $externalIPListIsEmpty) (not $externalIPListEqualsReplicaCount)) (and $externalIPListIsEmpty $nodePortListIsEmpty)) -}}
kafka: .Values.externalAccess.controller.service.externalIPs
    Number of controller-eligible replicas and externalAccess.controller.service.externalIPs array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length externalIPs = {{ $externalIPListLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as externalIPs list */}}
{{- define "kafka.validateValues.broker.externalIPListLength" -}}
{{- $replicaCount := int .Values.broker.replicaCount -}}
{{- $externalIPListLength := len .Values.externalAccess.broker.service.externalIPs -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.broker.service.externalIPs -}}
{{- $externalIPListEqualsReplicaCount := eq $externalIPListLength $replicaCount -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.broker.service.nodePorts -}}
{{- if and .Values.externalAccess.enabled (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq .Values.externalAccess.broker.service.type "NodePort") (or (and (not $externalIPListIsEmpty) (not $externalIPListEqualsReplicaCount)) (and $externalIPListIsEmpty $nodePortListIsEmpty)) -}}
kafka: .Values.externalAccess.broker.service.externalIPs
    Number of broker replicas and externalAccess.broker.service.externalIPs array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length externalIPs = {{ $externalIPListLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - broker rack assignment allowed values */}}
{{- define "kafka.validateValues.brokerRackAwareness" -}}
{{- if and .Values.brokerRackAwareness.enabled (ne .Values.brokerRackAwareness.cloudProvider "aws-az") (ne .Values.brokerRackAwareness.cloudProvider "azure") -}}
kafka: .Values.brokerRackAwareness.cloudProvider
    Available values for the cloud provider to use for broker rack awareness are "aws-az" or "azure"
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - domain must be defined if external service type ClusterIP */}}
{{- define "kafka.validateValues.domainSpecified" -}}
{{- if and (eq .Values.externalAccess.controller.service.type "ClusterIP") (empty .Values.externalAccess.controller.service.domain) -}}
kafka: .Values.externalAccess.controller.service.domain
    Domain must be specified if service type ClusterIP is set for external service
{{- end -}}
{{- if and (eq .Values.externalAccess.broker.service.type "ClusterIP") (empty .Values.externalAccess.broker.service.domain) -}}
kafka: .Values.externalAccess.broker.service.domain
    Domain must be specified if service type ClusterIP is set for external service
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - service type for external access */}}
{{- define "kafka.validateValues.externalAccessServiceType" -}}
{{- if and (not (eq .Values.externalAccess.controller.service.type "NodePort")) (not (eq .Values.externalAccess.controller.service.type "LoadBalancer")) (not (eq .Values.externalAccess.controller.service.type "ClusterIP")) -}}
kafka: externalAccess.controller.service.type
    Available service type for external access are NodePort, LoadBalancer or ClusterIP.
{{- end -}}
{{- if and (not (eq .Values.externalAccess.broker.service.type "NodePort")) (not (eq .Values.externalAccess.broker.service.type "LoadBalancer")) (not (eq .Values.externalAccess.broker.service.type "ClusterIP")) -}}
kafka: externalAccess.broker.service.type
    Available service type for external access are NodePort, LoadBalancer or ClusterIP.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - RBAC should be enabled when autoDiscovery is enabled */}}
{{- define "kafka.validateValues.externalAccessAutoDiscoveryRBAC" -}}
{{- if and .Values.externalAccess.enabled .Values.defaultInitContainers.autoDiscovery.enabled (not .Values.rbac.create ) }}
kafka: rbac.create
    By specifying "externalAccess.enabled=true" and "defaultInitContainers.autoDiscovery.enabled=true"
    an initContainer will be used to auto-detect the external IPs/ports by querying the
    K8s API. Please note this initContainer requires specific RBAC resources. You can create them
    by specifying "--set rbac.create=true".
{{- end -}}
{{- if and .Values.externalAccess.enabled .Values.defaultInitContainers.autoDiscovery.enabled (gt (int .Values.controller.replicaCount) 0) (not .Values.controller.automountServiceAccountToken) }}
kafka: controller-automountServiceAccountToken
    By specifying "externalAccess.enabled=true" and "defaultInitContainers.autoDiscovery.enabled=true"
    an initContainer will be used to auto-detect the external IPs/ports by querying the
    K8s API. Please note this initContainer requires the service account token. Please set controller.automountServiceAccountToken=true
    and broker.automountServiceAccountToken=true.
{{- end -}}
{{- if and .Values.externalAccess.enabled .Values.defaultInitContainers.autoDiscovery.enabled (gt (int .Values.broker.replicaCount) 0) (not .Values.broker.automountServiceAccountToken) }}
kafka: broker-automountServiceAccountToken
    By specifying "externalAccess.enabled=true" and "defaultInitContainers.autoDiscovery.enabled=true"
    an initContainer will be used to auto-detect the external IPs/ports by querying the
    K8s API. Please note this initContainer requires the service account token. Please set controller.automountServiceAccountToken=true
    and broker.automountServiceAccountToken=true.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - LoadBalancerIPs or LoadBalancerNames should be set when autoDiscovery is disabled */}}
{{- define "kafka.validateValues.externalAccessAutoDiscoveryIPsOrNames" -}}
{{- $loadBalancerNameListLength := len .Values.externalAccess.controller.service.loadBalancerNames -}}
{{- $loadBalancerIPListLength := len .Values.externalAccess.controller.service.loadBalancerIPs -}}
{{- if and .Values.externalAccess.enabled (gt (int .Values.controller.replicaCount) 0) (or .Values.externalAccess.controller.forceExpose (not .Values.controller.controllerOnly)) (eq .Values.externalAccess.controller.service.type "LoadBalancer") (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq $loadBalancerNameListLength 0) (eq $loadBalancerIPListLength 0) }}
kafka: externalAccess.controller.service.loadBalancerNames or externalAccess.controller.service.loadBalancerIPs
    By specifying "externalAccess.enabled=true", "defaultInitContainers.autoDiscovery.enabled=false" and
    "externalAccess.controller.service.type=LoadBalancer" at least one of externalAccess.controller.service.loadBalancerNames
    or externalAccess.controller.service.loadBalancerIPs must be set and the length of those arrays must be equal
    to the number of replicas.
{{- end -}}
{{- $loadBalancerNameListLength := len .Values.externalAccess.broker.service.loadBalancerNames -}}
{{- $loadBalancerIPListLength := len .Values.externalAccess.broker.service.loadBalancerIPs -}}
{{- $replicaCount := int .Values.broker.replicaCount }}
{{- if and .Values.externalAccess.enabled (gt 0 $replicaCount) (eq .Values.externalAccess.broker.service.type "LoadBalancer") (not .Values.defaultInitContainers.autoDiscovery.enabled) (eq $loadBalancerNameListLength 0) (eq $loadBalancerIPListLength 0) }}
kafka: externalAccess.broker.service.loadBalancerNames or externalAccess.broker.service.loadBalancerIPs
    By specifying "externalAccess.enabled=true", "defaultInitContainers.autoDiscovery.enabled=false" and
    "externalAccess.broker.service.type=LoadBalancer" at least one of externalAccess.broker.service.loadBalancerNames
    or externalAccess.broker.service.loadBalancerIPs must be set and the length of those arrays must be equal
    to the number of replicas.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as loadBalancerIPs list */}}
{{- define "kafka.validateValues.externalAccessServiceList" -}}
{{- $replicaCount := int .context.Values.controller.replicaCount }}
{{- $listLength := len (get .context.Values.externalAccess.controller.service .element) -}}
{{- if and .context.Values.externalAccess.enabled (or .context.Values.externalAccess.controller.forceExpose (not .context.Values.controller.controllerOnly)) (not .context.Values.defaultInitContainers.autoDiscovery.enabled) (eq .context.Values.externalAccess.controller.service.type "LoadBalancer") (gt $listLength 0) (not (eq $replicaCount $listLength)) }}
kafka: externalAccess.service.{{ .element }}
    Number of replicas and {{ .element }} array length must be the same. Currently: replicaCount = {{ $replicaCount }} and {{ .element }} = {{ $listLength }}
{{- end -}}
{{- $replicaCount := int .context.Values.broker.replicaCount }}
{{- $listLength := len (get .context.Values.externalAccess.broker.service .element) -}}
{{- if and .context.Values.externalAccess.enabled (gt 0 $replicaCount) (not .context.Values.defaultInitContainers.autoDiscovery.enabled) (eq .context.Values.externalAccess.broker.service.type "LoadBalancer") (gt $listLength 0) (not (eq $replicaCount $listLength)) }}
kafka: externalAccess.service.{{ .element }}
    Number of replicas and {{ .element }} array length must be the same. Currently: replicaCount = {{ $replicaCount }} and {{ .element }} = {{ $listLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - SASL mechanisms must be provided when using SASL */}}
{{- define "kafka.validateValues.saslMechanisms" -}}
{{- if and (include "kafka.saslEnabled" .) (not .Values.sasl.enabledMechanisms) }}
kafka: sasl.enabledMechanisms
    The SASL mechanisms are required when listeners use SASL security protocol.
{{- end }}
{{- if not (contains .Values.sasl.interBrokerMechanism .Values.sasl.enabledMechanisms) }}
kafka: sasl.enabledMechanisms
    sasl.interBrokerMechanism must be provided and it should be one of the specified mechanisms at sasl.enabledMechanisms
{{- end -}}
{{- if not (contains .Values.sasl.controllerMechanism .Values.sasl.enabledMechanisms) }}
kafka: sasl.enabledMechanisms
    sasl.controllerMechanism must be provided and it should be one of the specified mechanisms at sasl.enabledMechanisms
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - Secrets containing TLS certs must be provided when TLS authentication is enabled */}}
{{- define "kafka.validateValues.tlsSecret" -}}
{{- if and (include "kafka.sslEnabled" .) (eq (upper .Values.tls.type) "JKS") (empty .Values.tls.existingSecret) (not .Values.tls.autoGenerated.enabled) }}
kafka: tls.existingSecret
    A secret containing the Kafka JKS keystores and truststore is required
    when TLS encryption in enabled and TLS format is "JKS"
{{- else if and (include "kafka.sslEnabled" .) (eq (upper .Values.tls.type) "PEM") (empty .Values.tls.existingSecret) (not .Values.tls.autoGenerated.enabled) }}
kafka: tls.existingSecret
    A secret containing the Kafka TLS certificates and keys is required
    when TLS encryption in enabled and TLS format is "PEM"
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka provisioning - keyPasswordSecretKey, keystorePasswordSecretKey or truststorePasswordSecretKey must not be used without passwordsSecret */}}
{{- define "kafka.validateValues.provisioning.tlsPasswords" -}}
{{- if and (regexFind "SSL" (upper .Values.listeners.client.protocol)) .Values.provisioning.enabled (not .Values.provisioning.auth.tls.passwordsSecret) }}
{{- if or .Values.provisioning.auth.tls.keyPasswordSecretKey .Values.provisioning.auth.tls.keystorePasswordSecretKey .Values.provisioning.auth.tls.truststorePasswordSecretKey }}
kafka: tls.keyPasswordSecretKey,tls.keystorePasswordSecretKey,tls.truststorePasswordSecretKey
    tls.keyPasswordSecretKey,tls.keystorePasswordSecretKey,tls.truststorePasswordSecretKey
    must not be used without passwordsSecret setted.
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - At least 1 controller is configured or controller.quorum.bootstrap.servers is set  */}}
{{- define "kafka.validateValues.missingController" -}}
{{- if and (le (int .Values.controller.replicaCount) 0) (not .Values.controller.quorumBootstrapServers) }}
kafka: Missing controller-eligible nodes
    No controller-eligible nodes have been configured.
{{- end -}}
{{- end -}}
