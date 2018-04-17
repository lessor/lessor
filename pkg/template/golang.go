package template

import (
	"bytes"
	"html/template"

	"github.com/pkg/errors"
)

// RenderGolang renders templateText with the params map as the top-level
// template parameter.
//
// For example, namespace should be defined `{{ index . "namespace" }}`
// See the docs for more info: https://golang.org/pkg/text/template/
func RenderGolang(templateText string, params map[string]string) (string, error) {
	tmpl, err := template.New("golang-template").Parse(templateText)
	if err != nil {
		return "", errors.Wrap(err, "error creating go template from input")
	}

	b := new(bytes.Buffer)
	if err := tmpl.Execute(b, params); err != nil {
		return "", errors.Wrap(err, "error executing go template with params")
	}

	return b.String(), nil
}
