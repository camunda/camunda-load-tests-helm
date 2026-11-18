{{/*
Expand the name of the chart.
*/}}
{{- define "camunda-load-tests.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Default OIDC settings.

These live in the template, not in values.yaml, so that a key present in
`oidc` or the legacy `saas` always means the user set it explicitly. If the
defaults were in values.yaml, `oidc.enabled: false` would silently override a
legacy `saas.enabled: true`.
*/}}
{{- define "camunda-load-tests.oidcDefaults" -}}
{{- $credentials := dict "existingSecret" "" "clientId" "" "clientSecret" "" "zeebeRestAddress" "" "authServer" "https://login.cloud.ultrawombat.com/oauth/token" "authType" "OAUTH" "zeebeGrpcAddress" "" "authorizationAudience" "zeebe.ultrawombat.com" -}}
{{- dict "enabled" false "credentials" $credentials | toYaml -}}
{{- end -}}

{{/*
Resolve the effective OIDC configuration.

Defaults are overridden by the deprecated `saas` key, which is in turn
overridden by `oidc`. mergeOverwrite is used so an explicit `oidc.enabled:
false` beats a legacy `saas.enabled: true`. The result is rendered as YAML and
parsed back with fromYaml at each call site.
*/}}
{{- define "camunda-load-tests.oidc" -}}
{{- $defaults := include "camunda-load-tests.oidcDefaults" . | fromYaml -}}
{{- mergeOverwrite $defaults (deepCopy (.Values.saas | default dict)) (deepCopy (.Values.oidc | default dict)) | toYaml -}}
{{- end -}}

{{/*
Generate the common annotations for deployments.
*/}}
{{- define "camunda-load-tests.annotations" -}}
{{- $oidc := include "camunda-load-tests.oidc" . | fromYaml -}}
{{ if $.Values.global.extraConfig -}}
checksum/config: {{ $.Values.global.extraConfig | toYaml | sha256sum }}
{{- end }}
{{ if $oidc.enabled -}}
checksum/saas-credentials: {{ include (print $.Template.BasePath "/credentials.yaml") . | sha256sum }}
{{- end }}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "camunda-load-tests.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create a default fully qualified credentials name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "camunda-load-tests.credentials-name" -}}
{{- $oidc := include "camunda-load-tests.oidc" . | fromYaml -}}
{{- $credentials := $oidc.credentials | default dict -}}
{{- if $credentials.existingSecret }}
{{- $credentials.existingSecret }}
{{- else }}
{{- printf "%s-credentials" .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "camunda-load-tests.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels for all manifests
*/}}
{{- define "camunda-load-tests.labels" -}}
helm.sh/chart: {{ include "camunda-load-tests.chart" . }}
{{ include "camunda-load-tests.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.global.commonLabels }}
{{ tpl (toYaml .Values.global.commonLabels) $ }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "camunda-load-tests.selectorLabels" -}}
app.kubernetes.io/name: {{ include "camunda-load-tests.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "camunda-load-tests.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "camunda-load-tests.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}