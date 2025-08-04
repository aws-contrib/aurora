{{- if not .Pending -}}
{{- println "No migration files to execute" }}
{{- else -}}
{{- println .Header }}
{{- range $i, $f := .Applied }}
	{{- println }}
	{{- $checkFailed := false }}
	{{- range $cf := $f.Checks }}
		{{- println " " (yellow "--") "checks before migrating version" (cyan $f.File.Version) }}
		{{- range $s := $cf.Stmts }}
			{{- if $s.Error }}
				{{- println "   " (red "->") (indent_ln $s.Stmt 7) }}
			{{- else }}
				{{- println "   " (cyan "->") (indent_ln $s.Stmt 7) }}
			{{- end }}
		{{- end }}
		{{- with $cf.Error }}
			{{- $checkFailed = true }}
			{{- println "   " (redBgWhiteFg .Text) }}
		{{- else }}
			{{- printf "  %s ok (%s)\n\n" (yellow "--") (yellow ($cf.End.Sub $cf.Start).String) }}
		{{- end }}
	{{- end }}
	{{- if $checkFailed }}
		{{- continue }} {{- /* No statements were applied. */}}
	{{- end }}
	{{- println " " (yellow "--") "migrating version" (cyan $f.File.Version) }}
	{{- range $f.Applied }}
		{{- println "   " (cyan "->") (indent_ln . 7) }}
	{{- end }}
	{{- with .Error }}
		{{- println "   " (redBgWhiteFg .Text) }}
	{{- else }}
		{{- printf "  %s ok (%s)\n" (yellow "--") (yellow (.End.Sub .Start).String) }}
	{{- end }}
{{- else }}
	{{- println }}
	{{- with .Error }}
		{{- println "   " (redBgWhiteFg .) }}
	{{- end }}
{{- end }}
{{- println }}
{{- println " " (cyan "-------------------------") }}
{{- println " " (.Summary "  ") }}
{{- end -}}
