package template_test

import (
	"testing"

	"github.com/lessor/lessor/pkg/template"
	"github.com/stretchr/testify/require"
)

func TestRenderGolang(t *testing.T) {
	const templateText = `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: kuard
  name: kuard
  namespace: {{ index . "namespace" }}
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
      - image: {{ index . "image" }}
        imagePullPolicy: IfNotPresent
        name: kuard`

	params := map[string]string{
		"namespace": "acme-labs",
		"image":     "gcr.io/kuar-demo/kuard-amd64:1",
	}

	actual, err := template.RenderGolang(templateText, params)
	require.NoError(t, err)
	require.Equal(t, rendered(), actual)
}
