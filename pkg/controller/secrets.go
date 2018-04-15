package controller

import (
	"strings"

	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// rehydrateSecrets copies secrets from a source namespace into a destination
// namespace
func (c *Controller) rehydrateSecrets(source, destination string) error {
	// list the secrets in the source namespace
	sourceSecrets, err := c.secretsLister.Secrets(source).List(labels.Everything())
	if err != nil {
		return errors.Wrapf(err, "error listing secrets in the source namespace: %s", source)
	}

	// list the secrets in the destination namespace
	destinationSecrets, err := c.secretsLister.Secrets(destination).List(labels.Everything())
	if err != nil {
		return errors.Wrapf(err, "error listing secrets in the destination namespace: %s", destination)
	}

	// create a constant-time lookup of secrets in the destination namespace
	destinationCache := map[string]*corev1.Secret{}
	for _, s := range destinationSecrets {
		destinationCache[s.Name] = s
	}

	// make sure every secret in the source namespace is up-to-date in the
	// destination namespace
	for _, s := range sourceSecrets {
		if strings.HasPrefix(s.Name, "default-token-") {
			continue
		}

		d, ok := destinationCache[s.Name]
		if !ok {
			level.Debug(c.logger).Log(
				"msg", "creating secret in namespace",
				"secret", s.Name,
				"namespace", destination,
			)

			newSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      s.Name,
					Namespace: destination,
				},
				Data: s.Data,
			}
			if _, err := c.kubeClient.CoreV1().Secrets(destination).Create(newSecret); err != nil {
				return errors.Wrapf(err, "error creating secret %s in namespace %s", newSecret.Name, destination)
			}

			continue
		}

		level.Debug(c.logger).Log(
			"msg", "updating secret in namespace",
			"secret", d.Name,
			"namespace", destination,
		)

		d.Data = s.Data
		if _, err := c.kubeClient.CoreV1().Secrets(destination).Update(d); err != nil {
			return errors.Wrapf(err, "error updating secret %s in namespace %s", d.Name, destination)
		}
	}

	return nil
}
