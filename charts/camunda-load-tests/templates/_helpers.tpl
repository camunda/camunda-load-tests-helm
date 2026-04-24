{{/*
Expand the name of the chart.
*/}}
{{- define "camunda-load-tests.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Generate the common annotations for deployments.
*/}}
{{- define "camunda-load-tests.annotations" -}}
{{- $annotations := list -}}
{{- if $.Values.global.extraConfig -}}
{{- $annotations = append $annotations (printf "checksum/config: %s" ($.Values.global.extraConfig | toYaml | sha256sum)) -}}
{{- end -}}
{{- if $.Values.saas.enabled -}}
{{- $annotations = append $annotations (printf "checksum/saas-credentials: %s" (include (print $.Template.BasePath "/credentials.yaml") . | sha256sum)) -}}
{{- end -}}
{{- join "\n" $annotations -}}
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
{{- if .Values.saas.credentials.existingSecret }}
{{- .Values.saas.credentials.existingSecret }}
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
Common labels
*/}}
{{- define "camunda-load-tests.labels" -}}
helm.sh/chart: {{ include "camunda-load-tests.chart" . }}
{{ include "camunda-load-tests.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
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
