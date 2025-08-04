Migration Status:
{{- if eq .Status "OK" }} {{ green .Status }}{{ end }}
{{- if eq .Status "PENDING" }} {{ yellow .Status }}{{ end }}
  {{ yellow "--" }} Current Version: {{ cyan .Current }}
{{- if gt .Total 0 }}{{ printf " (%s statements applied)" (yellow "%d" .Count) }}{{ end }}
  {{ yellow "--" }} Next Version:    {{ if .Next }}{{ cyan .Next }}{{ if .FromCheckpoint }} (checkpoint){{ end }}{{ else }}UNKNOWN{{ end }}
{{- if gt .Total 0 }}{{ printf " (%s statements left)" (yellow "%d" .Left) }}{{ end }}
  {{ yellow "--" }} Executed Files:  {{ len .Applied }}{{ if gt .Total 0 }} (last one partially){{ end }}
  {{ yellow "--" }} Pending Files:   {{ add (len .Pending) (len .OutOfOrder) }}{{ if .OutOfOrder }} ({{ if .Pending }}{{ len .OutOfOrder }} {{ end }}out of order){{ end }}
{{- if gt .Total 0 }}

Last migration attempt had errors:
  {{ yellow "--" }} SQL:   {{ .SQL }}
  {{ yellow "--" }} {{ red "ERROR:" }} {{ .Error }}
{{- else if and .OutOfOrder .Error }}

  {{ red "ERROR:" }} {{ .Error }}
{{- end }} {{ yellow "--" }} Next Version:    {{ if .Next }}{{ cyan .Next }}{{ else }}UNKNOWN{{ end }}
