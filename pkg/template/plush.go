package template

import (
	"github.com/gobuffalo/plush"
)

// RenderPlush renders templateText with the keys and values in params as
// top-level plush context variables.
//
// For example, namespace should be defined `<%= namespace %>`
// See thee docs for more info: https://github.com/gobuffalo/plush
func RenderPlush(templateText string, params map[string]string) (string, error) {
	ctx := plush.NewContext()

	for k, v := range params {
		ctx.Set(k, v)
	}

	return plush.Render(templateText, ctx)
}
