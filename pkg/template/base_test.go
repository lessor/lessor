package template_test

func rendered() string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    run: kuard
  name: kuard
  namespace: acme-labs
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
      - image: gcr.io/kuar-demo/kuard-amd64:1
        imagePullPolicy: IfNotPresent
        name: kuard`
}
