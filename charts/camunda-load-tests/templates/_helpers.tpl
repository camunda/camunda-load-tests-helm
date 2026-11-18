{{/*
Expand the name of the chart.
*/}}
{{- define "camunda-load-tests.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Strip nil entries recursively from a map so YAML `null` overrides (e.g.
`oidc: { enabled: null }` or `oidc: { credentials: null }`) fall through to
defaults or legacy `saas` values instead of clobbering them in mergeOverwrite.
*/}}
{{- define "camunda-load-tests.pruneNilMap" -}}
{{- $in := . -}}
{{- range $k, $v := $in -}}
  {{- if eq $v nil -}}
    {{- $_ := unset $in $k -}}
  {{- else if kindIs "map" $v -}}
    {{- include "camunda-load-tests.pruneNilMap" $v -}}
  {{- end -}}
{{- end -}}
{{- end -}}

{{/*
Resolve the effective OIDC configuration.

Defaults are overridden by the deprecated `saas` key, which is in turn
overridden by `oidc`. Defaults live here rather than in values.yaml so that
an omitted `oidc.enabled` does not mask a legacy `saas.enabled: true`, while an
explicit `oidc.enabled: false` still beats `saas.enabled: true` via
mergeOverwrite. Credential fields are stringified before toYaml so numeric
`--set` values survive fromYaml without float64 scientific-notation conversion.
*/}}
{{- define "camunda-load-tests.oidc" -}}
{{- $defaults := dict
  "enabled" false
  "credentials" (dict
    "existingSecret" ""
    "clientId" ""
    "clientSecret" ""
    "zeebeRestAddress" ""
    "authServer" "https://login.cloud.ultrawombat.com/oauth/token"
    "authType" "OAUTH"
    "zeebeGrpcAddress" ""
    "authorizationAudience" "zeebe.ultrawombat.com"
  )
-}}
{{- $saas := deepCopy (.Values.saas | default dict) -}}
{{- include "camunda-load-tests.pruneNilMap" $saas -}}
{{- $oidc := deepCopy (.Values.oidc | default dict) -}}
{{- include "camunda-load-tests.pruneNilMap" $oidc -}}
{{- $merged := mergeOverwrite $defaults $saas $oidc -}}
{{- $creds := $merged.credentials | default dict -}}
{{- range $k, $v := $creds -}}
  {{- $_ := set $creds $k (toString $v) -}}
{{- end -}}
{{- $_ := set $merged "credentials" $creds -}}
{{- $merged | toYaml -}}
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