{{/*
This file contains helper templates for naming, labels, etc.
*/}}

{{- define "tiger-tail-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "tiger-tail-chart.fullname" -}}
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

{{- define "tiger-tail-chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tiger-tail-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "tiger-tail-chart.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{ include "tiger-tail-chart.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "tiger-tail-chart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "tiger-tail-chart.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Common pod security context
*/}}
{{- define "tiger-tail-chart.podSecurityContext" -}}
{{- if .Values.podSecurityContext.enabled }}
securityContext:
  runAsUser: {{ .Values.podSecurityContext.runAsUser }}
  runAsGroup: {{ .Values.podSecurityContext.runAsGroup }}
  fsGroup: {{ .Values.podSecurityContext.fsGroup }}
{{- end }}
{{- end }}

{{/*
Common container security context
*/}}
{{- define "tiger-tail-chart.securityContext" -}}
securityContext:
  {{- toYaml .Values.securityContext | nindent 12 }}
{{- end }}

{{/*
Common environment variables
*/}}
{{- define "tiger-tail-chart.commonEnv" -}}
# Basic env for DB and Redis usage:
- name: USE_REAL_DB
  value: {{ .Values.env.useRealDb | quote }}
- name: USE_REAL_REDIS
  value: {{ .Values.env.useRealRedis | quote }}
- name: SERVER_PORT
  value: {{ .Values.env.serverPort | quote }}
- name: SERVER_HOST
  value: {{ .Values.env.serverHost | quote }}
- name: LOG_LEVEL
  value: {{ .Values.env.logLevel | quote }}

# DB host references subchart if enabled; fallback to localhost if disabled:
- name: DB_HOST
  value: "{{ if .Values.postgresql.enabled }}{{ include "tiger-tail-chart.fullname" . }}-postgresql{{ else }}localhost{{ end }}"
- name: DB_PORT
  value: "5432"

# Pull DB credentials from our K8s Secret:
- name: DB_USER
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: db-user
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: db-password
- name: DB_NAME
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: db-name

# Redis host references subchart if enabled; fallback to localhost if disabled:
- name: REDIS_HOST
  value: "{{ if .Values.redis.enabled }}{{ include "tiger-tail-chart.fullname" . }}-redis-master{{ else }}localhost{{ end }}"
- name: REDIS_PORT
  value: "6379"
- name: REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: redis-password

# Basic Auth credentials:
- name: AUTH_USERNAME
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: auth-username
- name: AUTH_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "tiger-tail-chart.fullname" . }}-secrets
      key: auth-password
{{- end }}

{{/*
Helper to ensure subchart values match secrets
This is used in templates that need to pass values to subcharts
*/}}
{{- define "tiger-tail-chart.subchartValues" -}}
{{- if .Values.postgresql.enabled }}
postgresql:
  auth:
    username: {{ .Values.secrets.dbUser }}
    password: {{ .Values.secrets.dbPassword }}
    database: {{ .Values.secrets.dbName }}
{{- end }}
{{- if .Values.redis.enabled }}
redis:
  auth:
    password: {{ .Values.secrets.redisPassword }}
{{- end }}
{{- end }}

{{/*
Common container probes
*/}}
{{- define "tiger-tail-chart.probes" -}}
livenessProbe:
  httpGet:
    path: /livez
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /readyz
    port: http
  initialDelaySeconds: 5
  periodSeconds: 10
{{- end }}
