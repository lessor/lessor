package controller

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	clientset "github.com/lessor/lessor/pkg/client/clientset/versioned"
	lessorScheme "github.com/lessor/lessor/pkg/client/clientset/versioned/scheme"
	informers "github.com/lessor/lessor/pkg/client/informers/externalversions"
	listers "github.com/lessor/lessor/pkg/client/listers/lessor/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslistersv1beta2 "k8s.io/client-go/listers/apps/v1beta2"
	corelisters "k8s.io/client-go/listers/core/v1"
	policylisters "k8s.io/client-go/listers/policy/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

// Controller is the controller implementation for managing Tenant resources
type Controller struct {
	logger log.Logger

	// kubeClient is a standard kubernetes clientset
	kubeClient kubernetes.Interface
	// lessorClient is a clientset for our own API group
	lessorClient clientset.Interface

	namespacesLister corelisters.NamespaceLister
	namespacesSynced cache.InformerSynced

	secretsLister corelisters.SecretLister
	secretsSynced cache.InformerSynced

	deploymentsLister appslistersv1beta2.DeploymentLister
	deploymentsSynced cache.InformerSynced

	statefulSetsLister appslistersv1beta2.StatefulSetLister
	statefulSetsSynced cache.InformerSynced

	servicesLister corelisters.ServiceLister
	servicesSynced cache.InformerSynced

	podDisruptionBudgetsLister policylisters.PodDisruptionBudgetLister
	podDisruptionBudgetsSynced cache.InformerSynced

	tenantsLister listers.TenantLister
	tenantsSynced cache.InformerSynced

	// these are rate limited work queues. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	tenantWorkqueue workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new controller
func NewController(
	logger log.Logger,
	kubeClient kubernetes.Interface,
	lessorClient clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	lessorInformerFactory informers.SharedInformerFactory,
	broadcastEvents bool,
) *Controller {
	// Add cloud-operator types to the default Kubernetes Scheme so Events can be
	// logged for cloud-operator types
	lessorScheme.AddToScheme(scheme.Scheme)

	// Create event broadcaster
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(func(format string, args ...interface{}) {
		if broadcastEvents {
			level.Info(logger).Log(
				"from", "event broadcaster",
				"msg", fmt.Sprintf(format, args...),
			)
		}
	})
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: kubeClient.CoreV1().Events(""),
		},
	)
	recorder := eventBroadcaster.NewRecorder(
		scheme.Scheme,
		corev1.EventSource{
			Component: ControllerAgentName,
		},
	)

	// Get references to shared index informers for the Deployment and Tenant types.
	namespaceInformer := kubeInformerFactory.Core().V1().Namespaces()
	secretInformer := kubeInformerFactory.Core().V1().Secrets()
	deploymentInformer := kubeInformerFactory.Apps().V1beta2().Deployments()
	statefullSetInformer := kubeInformerFactory.Apps().V1beta2().StatefulSets()
	serviceInformer := kubeInformerFactory.Core().V1().Services()
	podDisruptionBudgetInformer := kubeInformerFactory.Policy().V1beta1().PodDisruptionBudgets()

	tenantInformer := lessorInformerFactory.Lessor().V1().Tenants()

	controller := &Controller{
		logger:                     logger,
		kubeClient:                 kubeClient,
		lessorClient:               lessorClient,
		namespacesLister:           namespaceInformer.Lister(),
		namespacesSynced:           namespaceInformer.Informer().HasSynced,
		secretsLister:              secretInformer.Lister(),
		secretsSynced:              secretInformer.Informer().HasSynced,
		deploymentsLister:          deploymentInformer.Lister(),
		deploymentsSynced:          deploymentInformer.Informer().HasSynced,
		statefulSetsLister:         statefullSetInformer.Lister(),
		statefulSetsSynced:         statefullSetInformer.Informer().HasSynced,
		servicesLister:             serviceInformer.Lister(),
		servicesSynced:             serviceInformer.Informer().HasSynced,
		podDisruptionBudgetsLister: podDisruptionBudgetInformer.Lister(),
		podDisruptionBudgetsSynced: podDisruptionBudgetInformer.Informer().HasSynced,
		tenantsLister:              tenantInformer.Lister(),
		tenantsSynced:              tenantInformer.Informer().HasSynced,
		tenantWorkqueue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Tenants"),
		recorder:                   recorder,
	}

	// Set up an event handler for when tenant resources change
	tenantInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueTenant,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueTenant(new)
		},
	})
	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workerCount int, stop <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.tenantWorkqueue.ShutDown()

	level.Debug(c.logger).Log("msg", "starting lessor controller")

	// Wait for the caches to be synced before starting workers
	level.Debug(c.logger).Log("msg", "waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stop,
		c.namespacesSynced,
		c.secretsSynced,
		c.deploymentsSynced,
		c.statefulSetsSynced,
		c.servicesSynced,
		c.podDisruptionBudgetsSynced,
		c.tenantsSynced,
	); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	// Launch workers to process tenant resources
	level.Debug(c.logger).Log("msg", "starting tenant workers", "count", workerCount)
	for i := 0; i < workerCount; i++ {
		go wait.Until(c.runTenantWorker, time.Second, stop)
	}

	<-stop
	level.Info(c.logger).Log("msg", "shutting down workers")

	return nil
}
