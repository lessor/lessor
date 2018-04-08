package controller

import (
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (c *Controller) updateOrCreateNamespace(n *corev1.Namespace) (*corev1.Namespace, error) {
	existing, err := c.namespacesLister.Get(n.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.CoreV1().Namespaces().Create(n)
	case err != nil:
		return nil, errors.Wrap(err, "error getting a namespace")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, n)
		return c.kubeClient.CoreV1().Namespaces().Update(merged)
	}
}

func (c *Controller) updateOrCreateDeployment(d *appsv1beta2.Deployment) (*appsv1beta2.Deployment, error) {
	existing, err := c.deploymentsLister.Deployments(d.Namespace).Get(d.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.AppsV1beta2().Deployments(d.Namespace).Create(d)
	case err != nil:
		return nil, errors.Wrap(err, "error getting deployment")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, d)
		return c.kubeClient.AppsV1beta2().Deployments(d.Namespace).Update(merged)
	}
}

func (c *Controller) updateOrCreateSecret(s *corev1.Secret) (*corev1.Secret, error) {
	existing, err := c.secretsLister.Secrets(s.Namespace).Get(s.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.CoreV1().Secrets(s.Namespace).Create(s)
	case err != nil:
		return nil, errors.Wrap(err, "error getting secret")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, s)
		return c.kubeClient.CoreV1().Secrets(s.Namespace).Update(merged)
	}
}

func (c *Controller) updateOrCreateService(s *corev1.Service) (*corev1.Service, error) {
	existing, err := c.servicesLister.Services(s.Namespace).Get(s.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.CoreV1().Services(s.Namespace).Create(s)
	case err != nil:
		return nil, errors.Wrap(err, "error getting service")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, s)
		return c.kubeClient.CoreV1().Services(s.Namespace).Update(merged)
	}
}

func (c *Controller) updateOrCreateStatefulSet(ss *appsv1beta2.StatefulSet) (*appsv1beta2.StatefulSet, error) {
	existing, err := c.statefulSetsLister.StatefulSets(ss.Namespace).Get(ss.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.AppsV1beta2().StatefulSets(ss.Namespace).Create(ss)
	case err != nil:
		return nil, errors.Wrap(err, "error getting stateful set")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, ss)
		return c.kubeClient.AppsV1beta2().StatefulSets(ss.Namespace).Update(merged)
	}
}

func (c *Controller) updateOrCreatePodDisruptionBudget(pdb *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error) {
	existing, err := c.podDisruptionBudgetsLister.PodDisruptionBudgets(pdb.Namespace).Get(pdb.Name)
	switch {
	case apierrors.IsNotFound(err):
		return c.kubeClient.PolicyV1beta1().PodDisruptionBudgets(pdb.Namespace).Create(pdb)
	case err != nil:
		return nil, errors.Wrap(err, "error getting pod disruption budget")
	default:
		merged := existing.DeepCopy()
		mergo.Merge(&merged, pdb)
		return c.kubeClient.PolicyV1beta1().PodDisruptionBudgets(pdb.Namespace).Update(merged)
	}
}
