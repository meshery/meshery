{{/*
Copyright Broadcom, Inc. All Rights Reserved.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/* vim: set filetype=mustache: */}}

{{/*
Returns an init-container that changes the owner and group of the persistent volume(s) mountpoint(s) to 'runAsUser:fsGroup' on each node
*/}}
{{- define "kafka.defaultInitContainers.volumePermissions" -}}
{{- $roleValues := index .context.Values .role -}}
- name: volume-permissions
  image: {{ include "kafka.volumePermissions.image" .context }}
  imagePullPolicy: {{ .context.Values.defaultInitContainers.volumePermissions.image.pullPolicy | quote }}
  {{- if .context.Values.defaultInitContainers.volumePermissions.containerSecurityContext.enabled }}
  securityContext: {{- include "common.compatibility.renderSecurityContext" (dict "secContext" .context.Values.defaultInitContainers.volumePermissions.containerSecurityContext "context" .context) | nindent 4 }}
  {{- end }}
  {{- if .context.Values.defaultInitContainers.volumePermissions.resources }}
  resources: {{- toYaml .context.Values.defaultInitContainers.volumePermissions.resources | nindent 4 }}
  {{- else if ne .context.Values.defaultInitContainers.volumePermissions.resourcesPreset "none" }}
  resources: {{- include "common.resources.preset" (dict "type" .context.Values.defaultInitContainers.volumePermissions.resourcesPreset) | nindent 4 }}
  {{- end }}
  command:
    - /bin/bash
  args:
    - -ec
    - |
      mkdir -p {{ $roleValues.persistence.mountPath }} {{ $roleValues.logPersistence.mountPath }}
      {{- if eq ( toString ( .context.Values.defaultInitContainers.volumePermissions.containerSecurityContext.runAsUser )) "auto" }}
      find {{ $roleValues.persistence.mountPath }} -mindepth 1 -maxdepth 1 -not -name ".snapshot" -not -name "lost+found" |  xargs -r chown -R $(id -u):$(id -G | cut -d " " -f2)
      find {{ $roleValues.logPersistence.mountPath }} -mindepth 1 -maxdepth 1 -not -name ".snapshot" -not -name "lost+found" |  xargs -r chown -R $(id -u):$(id -G | cut -d " " -f2)
      {{- else }}
      find {{ $roleValues.persistence.mountPath }} -mindepth 1 -maxdepth 1 -not -name ".snapshot" -not -name "lost+found" |  xargs -r chown -R {{ $roleValues.containerSecurityContext.runAsUser }}:{{ $roleValues.podSecurityContext.fsGroup }}
      find {{ $roleValues.logPersistence.mountPath }} -mindepth 1 -maxdepth 1 -not -name ".snapshot" -not -name "lost+found" |  xargs -r chown -R {{ $roleValues.containerSecurityContext.runAsUser }}:{{ $roleValues.podSecurityContext.fsGroup }}
      {{- end }}
  volumeMounts:
    - name: data
      mountPath: {{ $roleValues.persistence.mountPath }}
    - name: logs
      mountPath: {{ $roleValues.logPersistence.mountPath }}
{{- end -}}

{{/*
Returns an init-container that auto-discovers the external access details
*/}}
{{- define "kafka.defaultInitContainers.autoDiscovery" -}}
{{- $externalAccess := index .context.Values.externalAccess .role }}
- name: auto-discovery
  image: {{ include "kafka.autoDiscovery.image" .context }}
  imagePullPolicy: {{ .context.Values.defaultInitContainers.autoDiscovery.image.pullPolicy | quote }}
  {{- if .context.Values.defaultInitContainers.autoDiscovery.containerSecurityContext.enabled }}
  securityContext: {{- include "common.compatibility.renderSecurityContext" (dict "secContext" .context.Values.defaultInitContainers.autoDiscovery.containerSecurityContext "context" .context) | nindent 4 }}
  {{- end }}
  {{- if .context.Values.defaultInitContainers.autoDiscovery.resources }}
  resources: {{- toYaml .context.Values.defaultInitContainers.autoDiscovery.resources | nindent 4 }}
  {{- else if ne .context.Values.defaultInitContainers.autoDiscovery.resourcesPreset "none" }}
  resources: {{- include "common.resources.preset" (dict "type" .context.Values.defaultInitContainers.autoDiscovery.resourcesPreset) | nindent 4 }}
  {{- end }}
  command:
    - /bin/bash
  args:
    - -ec
    - |
      SVC_NAME="${MY_POD_NAME}-external"
      AUTODISCOVERY_SERVICE_TYPE="${AUTODISCOVERY_SERVICE_TYPE:-}"

      # Auxiliary functions
      retry_while() {
          local -r cmd="${1:?cmd is missing}"
          local -r retries="${2:-12}"
          local -r sleep_time="${3:-5}"
          local return_value=1
          read -r -a command <<< "$cmd"
          for ((i = 1 ; i <= retries ; i+=1 )); do
              "${command[@]}" && return_value=0 && break
              sleep "$sleep_time"
          done
          return $return_value
      }
      k8s_svc_lb_ip() {
          local namespace=${1:?namespace is missing}
          local service=${2:?service is missing}
          local service_ip=$(kubectl get svc "$service" -n "$namespace" -o jsonpath="{.status.loadBalancer.ingress[0].ip}")
          local service_hostname=$(kubectl get svc "$service" -n "$namespace" -o jsonpath="{.status.loadBalancer.ingress[0].hostname}")
          if [[ -n ${service_ip} ]]; then
              echo "${service_ip}"
          else
              echo "${service_hostname}"
          fi
      }
      k8s_svc_lb_ip_ready() {
          local namespace=${1:?namespace is missing}
          local service=${2:?service is missing}
          [[ -n "$(k8s_svc_lb_ip "$namespace" "$service")" ]]
      }
      k8s_svc_node_port() {
          local namespace=${1:?namespace is missing}
          local service=${2:?service is missing}
          local index=${3:-0}
          local node_port="$(kubectl get svc "$service" -n "$namespace" -o jsonpath="{.spec.ports[$index].nodePort}")"
          echo "$node_port"
      }

      if [[ "$AUTODISCOVERY_SERVICE_TYPE" = "LoadBalancer" ]]; then
          # Wait until LoadBalancer IP is ready
          retry_while "k8s_svc_lb_ip_ready $MY_POD_NAMESPACE $SVC_NAME" || exit 1
          # Obtain LoadBalancer external IP
          k8s_svc_lb_ip "$MY_POD_NAMESPACE" "$SVC_NAME" | tee "/shared/external-host.txt"
      elif [[ "$AUTODISCOVERY_SERVICE_TYPE" = "NodePort" ]]; then
          k8s_svc_node_port "$MY_POD_NAMESPACE" "$SVC_NAME" | tee "/shared/external-port.txt"
      else
          echo "Unsupported autodiscovery service type: '$AUTODISCOVERY_SERVICE_TYPE'"
          exit 1
      fi

  env:
    - name: MY_POD_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
    - name: MY_POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    - name: AUTODISCOVERY_SERVICE_TYPE
      value: {{ $externalAccess.service.type | quote }}
  volumeMounts:
    - name: init-shared
      mountPath: /shared
{{- end -}}

{{/*
Returns an init-container that prepares the Kafka configuration files for main containers to use them
*/}}
{{- define "kafka.defaultInitContainers.prepareConfig" -}}
{{- $roleValues := index .context.Values .role -}}
{{- $externalAccessEnabled := or (and (eq .role "broker") .context.Values.externalAccess.enabled) (and (eq .role "controller") .context.Values.externalAccess.enabled (or .context.Values.externalAccess.controller.forceExpose (not .context.Values.controller.controllerOnly))) }}
- name: prepare-config
  image: {{ include "kafka.image" .context }}
  imagePullPolicy: {{ .context.Values.image.pullPolicy }}
  {{- if .context.Values.defaultInitContainers.prepareConfig.containerSecurityContext.enabled }}
  securityContext: {{- include "common.compatibility.renderSecurityContext" (dict "secContext" .context.Values.defaultInitContainers.prepareConfig.containerSecurityContext "context" .context) | nindent 4 }}
  {{- end }}
  {{- if .context.Values.defaultInitContainers.prepareConfig.resources }}
  resources: {{- toYaml .context.Values.defaultInitContainers.prepareConfig.resources | nindent 4 }}
  {{- else if ne .context.Values.defaultInitContainers.prepareConfig.resourcesPreset "none" }}
  resources: {{- include "common.resources.preset" (dict "type" .context.Values.defaultInitContainers.prepareConfig.resourcesPreset) | nindent 4 }}
  {{- end }}
  command:
    - /bin/bash
  args:
    - -ec
    - |
      . /opt/bitnami/scripts/libkafka.sh

      {{- if $externalAccessEnabled }}
      configure_external_access() {
          local host port
          # Configure external hostname
          if [[ -f "/shared/external-host.txt" ]]; then
              host=$(cat "/shared/external-host.txt")
          elif [[ -n "${EXTERNAL_ACCESS_HOST:-}" ]]; then
              host="$EXTERNAL_ACCESS_HOST"
          elif [[ -n "${EXTERNAL_ACCESS_HOSTS_LIST:-}" ]]; then
              read -r -a hosts <<< "$(tr ',' ' ' <<<"${EXTERNAL_ACCESS_HOSTS_LIST}")"
              host="${hosts[$POD_ID]}"
          elif is_boolean_yes "$EXTERNAL_ACCESS_HOST_USE_PUBLIC_IP"; then
              host=$(curl -s https://ipinfo.io/ip)
          else
              error "External access hostname not provided"
          fi
          # Configure external port
          if [[ -f "/shared/external-port.txt" ]]; then
              port=$(cat "/shared/external-port.txt")
          elif [[ -n "${EXTERNAL_ACCESS_PORT:-}" ]]; then
              port="$EXTERNAL_ACCESS_PORT"
              if is_boolean_yes "${EXTERNAL_ACCESS_PORT_AUTOINCREMENT:-}"; then
                  port="$((port + POD_ID))"
              fi
          elif [[ -n "${EXTERNAL_ACCESS_PORTS_LIST:-}" ]]; then
              read -r -a ports <<<"$(tr ',' ' ' <<<"${EXTERNAL_ACCESS_PORTS_LIST}")"
              port="${ports[$POD_ID]}"
          else
              error "External access port not provided"
          fi
          # Configure Kafka advertised listeners
          sed -i -E "s|^(advertised\.listeners=\S+)$|\1,${EXTERNAL_ACCESS_LISTENER_NAME}://${host}:${port}|" "$KAFKA_CONF_FILE"
      }
      {{- end }}
      {{- if include "kafka.sslEnabled" .context }}
      configure_kafka_tls() {
          # Remove previously existing keystores and certificates, if any
          rm -f /certs/kafka.keystore.jks /certs/kafka.truststore.jks
          rm -f /certs/tls.crt /certs/tls.key /certs/ca.crt
          find /certs -name "xx*" -exec rm {} \;
          if [[ "${KAFKA_TLS_TYPE}" = "PEM" ]]; then
              # Copy PEM certificate and key
              if [[ -f "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.crt" && "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.key" ]]; then
                  cp "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.crt" /certs/tls.crt
                  # Copy the PEM key ensuring the key used PEM format with PKCS#8
                  openssl pkcs8 -topk8 -nocrypt -passin pass:"${KAFKA_TLS_PEM_KEY_PASSWORD:-}" -in "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.key" > /certs/tls.key
              elif [[ -f /mounted-certs/tls.crt && -f /mounted-certs/tls.key ]]; then
                  cp "/mounted-certs/tls.crt" /certs/tls.crt
                  # Copy the PEM key ensuring the key used PEM format with PKCS#8
                  openssl pkcs8 -topk8 -passin pass:"${KAFKA_TLS_PEM_KEY_PASSWORD:-}" -nocrypt -in "/mounted-certs/tls.key" > /certs/tls.key
              else
                  error "PEM key and cert files not found"
              fi
      {{- if not .context.Values.tls.pemChainIncluded }}
              # Copy CA certificate
              if [[ -f /mounted-certs/ca.crt ]]; then
                  cp /mounted-certs/ca.crt /certs/ca.crt
              else
                  error "CA certificate file not found"
              fi
      {{- else }}
              # CA certificates are also included in the same certificate
              # All public certs will be included in the truststore
              cp /certs/tls.crt /certs/ca.crt
      {{- end }}
              # Create JKS keystore from PEM cert and key
              openssl pkcs12 -export -in "/certs/tls.crt" \
                  -passout pass:"$KAFKA_TLS_KEYSTORE_PASSWORD" \
                  -inkey "/certs/tls.key" \
                  -out "/certs/kafka.keystore.p12"
              keytool -importkeystore -srckeystore "/certs/kafka.keystore.p12" \
                  -srcstoretype PKCS12 \
                  -srcstorepass "$KAFKA_TLS_KEYSTORE_PASSWORD" \
                  -deststorepass "$KAFKA_TLS_KEYSTORE_PASSWORD" \
                  -destkeystore "/certs/kafka.keystore.jks" \
                  -noprompt
              # Create JKS truststore from CA cert
              keytool -keystore /certs/kafka.truststore.jks -alias CARoot -import -file /certs/ca.crt -storepass "$KAFKA_TLS_TRUSTSTORE_PASSWORD" -noprompt
              # Remove extra files
              rm -f "/certs/kafka.keystore.p12" "/certs/tls.crt" "/certs/tls.key" "/certs/ca.crt"
          elif [[ "$KAFKA_TLS_TYPE" = "JKS" ]]; then
              if [[ -f "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.keystore.jks" ]]; then
                  cp "/mounted-certs/kafka-${POD_ROLE}-${POD_ID}.keystore.jks" /certs/kafka.keystore.jks
              elif [[ -f "$KAFKA_TLS_KEYSTORE_FILE" ]]; then
                  cp "$KAFKA_TLS_KEYSTORE_FILE" /certs/kafka.keystore.jks
              else
                  error "Keystore file not found"
              fi
              if [[ -f "$KAFKA_TLS_TRUSTSTORE_FILE" ]]; then
                  cp "$KAFKA_TLS_TRUSTSTORE_FILE" /certs/kafka.truststore.jks
              else
                  error "Truststore file not found"
              fi
          else
              error "Invalid type $KAFKA_TLS_TYPE"
          fi
          # Configure TLS password settings in Kafka configuration
          [[ -n "${KAFKA_TLS_KEYSTORE_PASSWORD:-}" ]] && kafka_server_conf_set "ssl.keystore.password" "$KAFKA_TLS_KEYSTORE_PASSWORD"
          [[ -n "${KAFKA_TLS_TRUSTSTORE_PASSWORD:-}" ]] && kafka_server_conf_set "ssl.truststore.password" "$KAFKA_TLS_TRUSTSTORE_PASSWORD"
          [[ -n "${KAFKA_TLS_PEM_KEY_PASSWORD:-}" ]] && kafka_server_conf_set "ssl.key.password" "$KAFKA_TLS_PEM_KEY_PASSWORD"
          # Avoid errors caused by previous checks
          true
      }
      {{- end }}
      {{- if include "kafka.saslEnabled" .context }}
      configure_kafka_sasl() {
          # Replace placeholders with passwords
      {{- if regexFind "SASL" (upper .context.Values.listeners.interbroker.protocol) }}
        {{- if include "kafka.saslUserPasswordsEnabled" .context }}
          replace_in_file "$KAFKA_CONF_FILE" "interbroker-password-placeholder" "$KAFKA_INTER_BROKER_PASSWORD"
        {{- end }}
        {{- if include "kafka.saslClientSecretsEnabled" .context }}
          replace_in_file "$KAFKA_CONF_FILE" "interbroker-client-secret-placeholder" "$KAFKA_INTER_BROKER_CLIENT_SECRET"
        {{- end }}
      {{- end }}
      {{- if regexFind "SASL" (upper .context.Values.listeners.controller.protocol) }}
        {{- if include "kafka.saslUserPasswordsEnabled" .context }}
          replace_in_file "$KAFKA_CONF_FILE" "controller-password-placeholder" "$KAFKA_CONTROLLER_PASSWORD"
        {{- end }}
        {{- if include "kafka.saslClientSecretsEnabled" .context }}
          replace_in_file "$KAFKA_CONF_FILE" "controller-client-secret-placeholder" "$KAFKA_CONTROLLER_CLIENT_SECRET"
       {{- end }}
      {{- end }}
      {{- if include "kafka.client.saslEnabled" .context }}
          read -r -a passwords <<< "$(tr ',;' ' ' <<<"${KAFKA_CLIENT_PASSWORDS:-}")"
          for ((i = 0; i < ${#passwords[@]}; i++)); do
              replace_in_file "$KAFKA_CONF_FILE" "password-placeholder-${i}\"" "${passwords[i]}\""
          done
      {{- end }}
      }
      {{- end }}
      {{- if .context.Values.brokerRackAwareness.enabled }}
      configure_kafka_broker_rack() {
          local -r metadata_api_ip="169.254.169.254"
          local broker_rack=""
      {{- if eq .context.Values.brokerRackAwareness.cloudProvider "aws-az" }}
          echo "Obtaining broker.rack for aws-az rack assignment"
          ec2_metadata_token=$(curl -X PUT "http://${metadata_api_ip}/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60")
          broker_rack=$(curl -H "X-aws-ec2-metadata-token: $ec2_metadata_token" "http://${metadata_api_ip}/latest/meta-data/placement/availability-zone-id")
      {{- else if eq .context.Values.brokerRackAwareness.cloudProvider "azure" }}
          echo "Obtaining broker.rack for azure rack assignment"
          location=$(curl -s -H Metadata:true --noproxy "*" "http://${metadata_api_ip}/metadata/instance/compute/location?api-version={{ .context.Values.brokerRackAwareness.azureApiVersion }}&format=text")
          zone=$(curl -s -H Metadata:true --noproxy "*" "http://${metadata_api_ip}/metadata/instance/compute/zone?api-version={{ .context.Values.brokerRackAwareness.azureApiVersion }}&format=text")
          broker_rack="${location}-${zone}"
      {{- end }}
          kafka_server_conf_set "broker.rack" "$broker_rack"
      }
      {{- end }}
      {{- if and $externalAccessEnabled .context.Values.defaultInitContainers.autoDiscovery.enabled }}
      # Wait for autodiscovery to finish
      retry_while "test -f /shared/external-host.txt -o -f /shared/external-port.txt" || error "Timed out waiting for autodiscovery init-container"
      {{- end }}

      cp /configmaps/server.properties $KAFKA_CONF_FILE

      # Get pod ID and role, last and second last fields in the pod name respectively
      POD_ID="${MY_POD_NAME##*-}"
      POD_ROLE="${MY_POD_NAME%-*}"; POD_ROLE="${POD_ROLE##*-}"

      # Configure node.id
      ID=$((POD_ID + KAFKA_MIN_ID))
      [[ -f "/bitnami/kafka/data/meta.properties" ]] && ID="$(grep "node.id" /bitnami/kafka/data/meta.properties | awk -F '=' '{print $2}')"
      kafka_server_conf_set "node.id" "$ID"
      # Configure initial controllers
      if [[ "controller" =~ "$POD_ROLE" ]]; then
          INITIAL_CONTROLLERS=()
          for ((i = 0; i < {{ int .context.Values.controller.replicaCount }}; i++)); do
              var="KAFKA_CONTROLLER_${i}_DIR_ID"; DIR_ID="${!var}"
              [[ $i -eq $POD_ID ]] && [[ -f "/bitnami/kafka/data/meta.properties" ]] && DIR_ID="$(grep "directory.id" /bitnami/kafka/data/meta.properties | awk -F '=' '{print $2}')"
              INITIAL_CONTROLLERS+=("${i}@${KAFKA_FULLNAME}-${POD_ROLE}-${i}.${KAFKA_CONTROLLER_SVC_NAME}.${MY_POD_NAMESPACE}.svc.${CLUSTER_DOMAIN}:${KAFKA_CONTROLLER_PORT}:${DIR_ID}")
          done
          echo "${INITIAL_CONTROLLERS[*]}" | awk -v OFS=',' '{$1=$1}1' > /shared/initial-controllers.txt
      fi
      {{- if not .context.Values.listeners.advertisedListeners }}
      replace_in_file "$KAFKA_CONF_FILE" "advertised-address-placeholder" "${MY_POD_NAME}.${KAFKA_FULLNAME}-${POD_ROLE}-headless.${MY_POD_NAMESPACE}.svc.${CLUSTER_DOMAIN}"
      {{- if $externalAccessEnabled }}
      configure_external_access
      {{- end }}
      {{- end }}
      {{- if include "kafka.sslEnabled" .context }}
      configure_kafka_tls
      {{- end }}
      {{- if include "kafka.saslEnabled" .context }}
      sasl_env_vars=(
        KAFKA_CLIENT_PASSWORDS
        KAFKA_INTER_BROKER_PASSWORD
        KAFKA_INTER_BROKER_CLIENT_SECRET
        KAFKA_CONTROLLER_PASSWORD
        KAFKA_CONTROLLER_CLIENT_SECRET
      )
      for env_var in "${sasl_env_vars[@]}"; do
          file_env_var="${env_var}_FILE"
          if [[ -n "${!file_env_var:-}" ]]; then
              if [[ -r "${!file_env_var:-}" ]]; then
                  export "${env_var}=$(< "${!file_env_var}")"
                  unset "${file_env_var}"
              else
                  warn "Skipping export of '${env_var}'. '${!file_env_var:-}' is not readable."
              fi
          fi
      done
      configure_kafka_sasl
      {{- end }}
      {{- if .context.Values.brokerRackAwareness.enabled }}
      configure_kafka_broker_rack
      {{- end }}
      if [[ -f /secret-config/server-secret.properties ]]; then
          cat /secret-config/server-secret.properties >> $KAFKA_CONF_FILE
      fi

      {{- include "common.tplvalues.render" ( dict "value" .context.Values.defaultInitContainers.prepareConfig.extraInit "context" .context ) | nindent 6 }}
  env:
    - name: BITNAMI_DEBUG
      value: {{ ternary "true" "false" (or .context.Values.image.debug .context.Values.diagnosticMode.enabled) | quote }}
    - name: MY_POD_NAME
      valueFrom:
        fieldRef:
            fieldPath: metadata.name
    - name: MY_POD_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
    - name: KAFKA_FULLNAME
      value: {{ include "common.names.fullname" .context | quote }}
    - name: CLUSTER_DOMAIN
      value: {{ .context.Values.clusterDomain | quote }}
    - name: KAFKA_VOLUME_DIR
      value: {{ $roleValues.persistence.mountPath | quote }}
    - name: KAFKA_CONF_FILE
      value: /config/server.properties
    - name: KAFKA_MIN_ID
      value: {{ $roleValues.minId | quote }}
    - name: KAFKA_CONTROLLER_SVC_NAME
      value: {{ printf "%s-headless" (include "kafka.controller.fullname" .context) | trunc 63 | trimSuffix "-" }}
    - name: KAFKA_CONTROLLER_PORT
      value: {{ .context.Values.listeners.controller.containerPort | quote }}
    {{- $kraftSecret := include "kafka.kraftSecretName" .context }}
    {{- range $i := until (int .context.Values.controller.replicaCount) }}
    - name: KAFKA_CONTROLLER_{{ $i }}_DIR_ID
      valueFrom:
        secretKeyRef:
          name: {{ $kraftSecret }}
          key: controller-{{ $i }}-id
    {{- end }}
    {{- if $externalAccessEnabled }}
    - name: EXTERNAL_ACCESS_LISTENER_NAME
      value: {{ upper .context.Values.listeners.external.name | quote }}
    {{- $externalAccess := index .context.Values.externalAccess .role }}
    {{- if or (eq $externalAccess.service.type "LoadBalancer") (and $externalAccess.service.loadBalancerNames (eq $externalAccess.service.type "ClusterIP")) }}
    {{- if not .context.Values.defaultInitContainers.autoDiscovery.enabled }}
    - name: EXTERNAL_ACCESS_HOSTS_LIST
      value: {{ join "," (default $externalAccess.service.loadBalancerIPs $externalAccess.service.loadBalancerNames) | quote }}
    {{- end }}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote }}
    {{- else if eq $externalAccess.service.type "NodePort" }}
    {{- if $externalAccess.service.domain }}
    - name: EXTERNAL_ACCESS_HOST
      value: {{ $externalAccess.service.domain | quote }}
    {{- else if and $externalAccess.service.usePodIPs .context.Values.defaultInitContainers.autoDiscovery.enabled }}
    - name: MY_POD_IP
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    - name: EXTERNAL_ACCESS_HOST
      value: "$(MY_POD_IP)"
    {{- else if or $externalAccess.service.useHostIPs .context.Values.defaultInitContainers.autoDiscovery.enabled }}
    - name: HOST_IP
      valueFrom:
        fieldRef:
          fieldPath: status.hostIP
    - name: EXTERNAL_ACCESS_HOST
      value: "$(HOST_IP)"
    {{- else if and $externalAccess.service.externalIPs (not .context.Values.defaultInitContainers.autoDiscovery.enabled) }}
    - name: EXTERNAL_ACCESS_HOSTS_LIST
      value: {{ join "," $externalAccess.service.externalIPs }}
    {{- else }}
    - name: EXTERNAL_ACCESS_HOST_USE_PUBLIC_IP
      value: "true"
    {{- end }}
    {{- if not .context.Values.defaultInitContainers.autoDiscovery.enabled }}
    {{- if and $externalAccess.service.externalIPs (empty $externalAccess.service.nodePorts)}}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote }}
    {{- else }}
    - name: EXTERNAL_ACCESS_PORTS_LIST
      value: {{ join "," $externalAccess.service.nodePorts | quote }}
    {{- end }}
    {{- end }}
    {{- else if eq $externalAccess.service.type "ClusterIP" }}
    - name: EXTERNAL_ACCESS_HOST
      value: {{ $externalAccess.service.domain | quote }}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote}}
    - name: EXTERNAL_ACCESS_PORT_AUTOINCREMENT
      value: "true"
    {{- end }}
    {{- end }}
    {{- if include "kafka.saslEnabled" .context }}
    {{- include "kafka.saslEnv" .context | nindent 4 }}
    {{- end }}
    {{- if include "kafka.sslEnabled" .context }}
    - name: KAFKA_TLS_TYPE
      value: {{ ternary "PEM" "JKS" (or .context.Values.tls.autoGenerated.enabled (eq (upper .context.Values.tls.type) "PEM")) }}
    {{- if eq (upper .context.Values.tls.type) "JKS" }}
    - name: KAFKA_TLS_KEYSTORE_FILE
      value: {{ printf "/mounted-certs/%s" ( default "kafka.keystore.jks" .context.Values.tls.jksKeystoreKey) | quote }}
    - name: KAFKA_TLS_TRUSTSTORE_FILE
      value: {{ printf "/mounted-certs/%s" ( default "kafka.truststore.jks" .context.Values.tls.jksTruststoreKey) | quote }}
    {{- end }}
    - name: KAFKA_TLS_KEYSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.passwordsSecretKeystoreKey | quote }}
    - name: KAFKA_TLS_TRUSTSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.passwordsSecretTruststoreKey | quote }}
    {{- if and (not .context.Values.tls.autoGenerated.enabled) (or .context.Values.tls.keyPassword (and .context.Values.tls.passwordsSecret .context.Values.tls.passwordsSecretPemPasswordKey)) }}
    - name: KAFKA_TLS_PEM_KEY_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ default "key-password" .context.Values.tls.passwordsSecretPemPasswordKey | quote }}
    {{- end }}
    {{- end }}
  volumeMounts:
    - name: data
      mountPath: /bitnami/kafka
    - name: kafka-config
      mountPath: /config
    - name: kafka-configmaps
      mountPath: /configmaps
    - name: kafka-secret-config
      mountPath: /secret-config
    - name: tmp
      mountPath: /tmp
    - name: init-shared
      mountPath: /shared
    {{- if include "kafka.sslEnabled" .context }}
    - name: kafka-shared-certs
      mountPath: /certs
    {{- if and (include "kafka.sslEnabled" .context) (or .context.Values.tls.existingSecret .context.Values.tls.autoGenerated.enabled) }}
    - name: kafka-certs
      mountPath: /mounted-certs
      readOnly: true
    {{- end }}
    {{- end }}
    {{- if and .context.Values.usePasswordFiles (include "kafka.saslEnabled" .context) }}
    - name: kafka-sasl
      mountPath: /opt/bitnami/kafka/config/secrets
      readOnly: true
    {{- end }}
{{- end -}}
