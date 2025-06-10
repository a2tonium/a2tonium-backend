package json_generator

const jsonTemplate = `{
  "name": "` + "`" + `{{.Name}}` + "`" + ` Course Certificate",
  "description": "Certificate of completion for the ` + "`" + "{{.Name}}" + "`" + " course." +
	` Awarded to: ` + "`" + "{{.IIN}}" + "`" + "." +
	`This NFT certifies successful completion of all modules.",
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
    {{- end }}]
}
`
