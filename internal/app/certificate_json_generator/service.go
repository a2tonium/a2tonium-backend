package certificate_json_generator

import (
	"bytes"
	"fmt"
	"text/template"
)

const jsonTemplate = `{
  "name": "{{.Name}}",
  "description": "{{.Description}}",
  "image": "{{.Image}}",
  "attributes": [
    {{- range $i, $attr := .Attributes }}
    {{- if $i}},{{end}}
    {
      "trait_type": "{{$attr.TraitType}}",
      {{- if isString $attr.Value }}
      "value": "{{$attr.Value}}"
      {{- else }}
      "value": {{$attr.Value}}
      {{- end }}
    }
    {{- end }}
  ],
  "quizGrades": [
    {{- range $i, $grade := .QuizGrades }}
      {{- if $i}}, {{end}}"{{$grade}}"
    {{- end }}
  ]
}
`

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

// GenerateCertificateJSON generates the JSON file from the Certificate data,
// saves it in the specified directory, and returns the absolute path to the file.
func GenerateCertificateJSON(cert Certificate) (string, error) {

	funcMap := template.FuncMap{
		"isString": isString,
	}

	tmpl, err := template.New("cert").Funcs(funcMap).Parse(jsonTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cert)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
