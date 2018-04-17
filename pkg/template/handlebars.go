package template

import (
	"github.com/aymerick/raymond"
)

// RenderHandlebars renders templateText with the params map as the top-level
// template parameter.
//
// For example, namespace should be defined `{{ namespace }}`
// See the docs for more info: https://github.com/aymerick/raymond
func RenderHandlebars(templateText string, params map[string]string) (string, error) {
	return raymond.Render(templateText, params)
}
