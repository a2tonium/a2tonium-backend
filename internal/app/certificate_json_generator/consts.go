package certificate_json_generator

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
