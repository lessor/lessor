package template_test

import (
	"testing"

	"github.com/lessor/lessor/pkg/template"
	"github.com/stretchr/testify/require"
)

func TestRenderHandlebars(t *testing.T) {
	const templateText = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: kuard
  name: kuard
  namespace: {{ namespace }}
spec:
  replicas: 1
  selector:
    matchLabels:
      run: kuard
  template:
    metadata:
      labels:
        run: kuard
    spec:
      containers:
      - image: {{ image }}
        imagePullPolicy: IfNotPresent
        name: kuard`

	params := map[string]string{
		"namespace": "acme-labs",
		"image":     "gcr.io/kuar-demo/kuard-amd64:1",
	}

	actual, err := template.RenderHandlebars(templateText, params)
	require.NoError(t, err)
	require.Equal(t, rendered(), actual)
}
