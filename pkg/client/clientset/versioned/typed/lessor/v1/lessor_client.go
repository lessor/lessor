package v1

import (
	v1 "github.com/lessor/lessor/pkg/apis/lessor.io/v1"
	"github.com/lessor/lessor/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type LessorV1Interface interface {
	RESTClient() rest.Interface
	TenantsGetter
}

// LessorV1Client is used to interact with features provided by the lessor.io group.
type LessorV1Client struct {
	restClient rest.Interface
}

func (c *LessorV1Client) Tenants(namespace string) TenantInterface {
	return newTenants(c, namespace)
}

// NewForConfig creates a new LessorV1Client for the given config.
func NewForConfig(c *rest.Config) (*LessorV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &LessorV1Client{client}, nil
}

// NewForConfigOrDie creates a new LessorV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *LessorV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new LessorV1Client for the given RESTClient.
func New(c rest.Interface) *LessorV1Client {
	return &LessorV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *LessorV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
