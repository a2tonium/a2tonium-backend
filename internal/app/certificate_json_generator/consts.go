package certificate_json_generator

import (
	"errors"
	"path/filepath"
)

var (
	templatesPattern = filepath.Join("internal", "app", "script_generator", "templates", "*.gohtml")
)

var (
	ErrInvalidVusNumber   = errors.New("INVALID_VUS_NUMBER")
	ErrInvalidDuration    = errors.New("INVALID_DURATION")
	ErrInvalidDestination = errors.New("INVALID_DESTINATION")
	ErrInvalidHttpMethod  = errors.New("INVALID_HTTP_METHOD")
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
  "content_url": "",
  "quizGrades": [
    {{- range $i, $grade := .QuizGrades }}
      {{- if $i}}, {{end}}"{{$grade}}"
    {{- end }}
  ]
}
`
