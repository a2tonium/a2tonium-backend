package certificate_json_generator

import (
	"bytes"
	"fmt"
	"html/template"
)

// GenerateCertificateJSON generates the JSON file from the Certificate data,
// saves it in the specified directory, and returns the absolute path to the file.
func (c *CertificateJsonGeneratorService) GenerateCertificateJSON(cert Certificate) (string, error) {
	funcMap := template.FuncMap{
		"isString": func(v interface{}) bool {
			_, ok := v.(string)
			return ok
		},
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
