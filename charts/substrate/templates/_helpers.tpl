{{/*
Copyright 2026 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/}}

{{/*
substrate.snapshotStorageEnv renders the S3-compatible storage env vars
every snapshot-plugin and atelet container needs. Rendered identically in
three places (ate-api-server's snapshot-plugin init container, atelet's
snapshot-plugin init container, and atelet's own main container), so it
lives here once instead of three times.

With rustfs.enabled (the default), it points at the bundled in-cluster
rustfs instance and its auto-generated credentials Secret. Otherwise it
requires snapshotStorage.* to name your own S3-compatible bucket and the
Secret holding its credentials - rustfs.yaml itself is skipped entirely in
that case (see rustfs.enabled's own comment in values.yaml).
*/}}
{{- define "substrate.snapshotStorageEnv" -}}
- name: ATE_STORAGE_BACKEND
  value: "s3"
{{- if .Values.rustfs.enabled }}
- name: AWS_ENDPOINT_URL
  value: "http://rustfs.ate-system.svc:9000"
- name: AWS_REGION
  value: "us-east-1"
- name: AWS_S3_USE_PATH_STYLE
  value: "true"
- name: AWS_ACCESS_KEY_ID
  valueFrom:
    secretKeyRef:
      name: rustfs-credentials
      key: access-key
- name: AWS_SECRET_ACCESS_KEY
  valueFrom:
    secretKeyRef:
      name: rustfs-credentials
      key: secret-key
{{- else }}
- name: AWS_ENDPOINT_URL
  value: {{ required "snapshotStorage.endpointURL is required when rustfs.enabled is false" .Values.snapshotStorage.endpointURL | quote }}
- name: AWS_REGION
  value: {{ required "snapshotStorage.region is required when rustfs.enabled is false" .Values.snapshotStorage.region | quote }}
- name: AWS_S3_USE_PATH_STYLE
  value: {{ .Values.snapshotStorage.usePathStyle | quote }}
- name: AWS_ACCESS_KEY_ID
  valueFrom:
    secretKeyRef:
      name: {{ required "snapshotStorage.credentialsSecretRef.name is required when rustfs.enabled is false" .Values.snapshotStorage.credentialsSecretRef.name | quote }}
      key: {{ .Values.snapshotStorage.credentialsSecretRef.accessKeyKey | quote }}
- name: AWS_SECRET_ACCESS_KEY
  valueFrom:
    secretKeyRef:
      name: {{ .Values.snapshotStorage.credentialsSecretRef.name | quote }}
      key: {{ .Values.snapshotStorage.credentialsSecretRef.secretKeyKey | quote }}
{{- end }}
{{- end -}}
