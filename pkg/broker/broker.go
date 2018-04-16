package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/pivotal-cf/brokerapi"
)

const (
	PlanNameShared    = "shared-vm"
	PlanNameDedicated = "dedicated-vm"
)

type InstanceCredentials struct {
	Host     string
	Port     int
	Password string
}

type InstanceCreator interface {
	Create(instanceID string) error
	Destroy(instanceID string) error
	InstanceExists(instanceID string) (bool, error)
}

type InstanceBinder interface {
	Bind(instanceID string, bindingID string) (InstanceCredentials, error)
	Unbind(instanceID string, bindingID string) error
	InstanceExists(instanceID string) (bool, error)
}

type Broker struct {
	InstanceCreators map[string]InstanceCreator
	InstanceBinders  map[string]InstanceBinder
}

func (b *Broker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	redisPlanList := []brokerapi.ServicePlan{}
	for _, plan := range b.redisPlans() {
		redisPlanList = append(redisPlanList, *plan)
	}

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          "redis01",
			Name:        "Redis",
			Description: "Redis is an in-memory data store.",
			Bindable:    true,
			Plans:       redisPlanList,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         "Redis",
				LongDescription:     "Redis is an in-memory data store.",
				DocumentationUrl:    "https://redis.io/documentation",
				SupportUrl:          "https://redis.io/support",
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", "https://redis.io/images/redis-small.png"),
				ProviderDisplayName: "Lessor",
			},
			Tags: []string{
				"redis",
			},
		},
	}, nil
}

func (b *Broker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}

	if b.instanceExists(instanceID) {
		return spec, brokerapi.ErrInstanceAlreadyExists
	}

	if serviceDetails.PlanID == "" {
		return spec, errors.New("plan_id required")
	}

	planIdentifier := ""
	for key, plan := range b.plans() {
		if plan.ID == serviceDetails.PlanID {
			planIdentifier = key
			break
		}
	}

	if planIdentifier == "" {
		return spec, errors.New("plan_id not recognized")
	}

	instanceCreator, ok := b.InstanceCreators[planIdentifier]
	if !ok {
		return spec, errors.New("instance creator not found for plan")
	}

	err = instanceCreator.Create(instanceID)
	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (b *Broker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	for _, instanceCreator := range b.InstanceCreators {
		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
		if instanceExists {
			return spec, instanceCreator.Destroy(instanceID)
		}
	}
	return spec, brokerapi.ErrInstanceDoesNotExist
}

func (b *Broker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}

	for _, repo := range b.InstanceBinders {
		instanceExists, _ := repo.InstanceExists(instanceID)
		if instanceExists {
			instanceCredentials, err := repo.Bind(instanceID, bindingID)
			if err != nil {
				return binding, err
			}
			credentialsMap := map[string]interface{}{
				"host":     instanceCredentials.Host,
				"port":     instanceCredentials.Port,
				"password": instanceCredentials.Password,
			}

			binding.Credentials = credentialsMap
			return binding, nil
		}
	}
	return brokerapi.Binding{}, brokerapi.ErrInstanceDoesNotExist
}

func (b *Broker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails) error {
	for _, repo := range b.InstanceBinders {
		instanceExists, _ := repo.InstanceExists(instanceID)
		if instanceExists {
			err := repo.Unbind(instanceID, bindingID)
			if err != nil {
				return brokerapi.ErrBindingDoesNotExist
			}
			return nil
		}
	}

	return brokerapi.ErrInstanceDoesNotExist
}

func (b *Broker) plans() map[string]*brokerapi.ServicePlan {
	combined := map[string]*brokerapi.ServicePlan{}

	for k, v := range b.redisPlans() {
		combined[k] = v
	}

	return combined
}

func (b *Broker) redisPlans() map[string]*brokerapi.ServicePlan {
	plans := map[string]*brokerapi.ServicePlan{}

	plans["redis-shared"] = &brokerapi.ServicePlan{
		ID:          "redis-shared-01",
		Name:        PlanNameShared,
		Description: "This plan provides a single Redis process on a shared VM, which is suitable for development and testing workloads",
		Metadata: &brokerapi.ServicePlanMetadata{
			Bullets: []string{
				"Each instance shares the same VM",
				"Single dedicated Redis process",
				"Suitable for development & testing workloads",
			},
			DisplayName: "Shared-VM",
		},
	}

	plans["redis-dedicated"] = &brokerapi.ServicePlan{
		ID:          "redis-dedicated-01",
		Name:        PlanNameDedicated,
		Description: "This plan provides a single Redis process on a dedicated VM, which is suitable for production workloads",
		Metadata: &brokerapi.ServicePlanMetadata{
			Bullets: []string{
				"Dedicated VM per instance",
				"Single dedicated Redis process",
				"Suitable for production workloads",
			},
			DisplayName: "Dedicated-VM",
		},
	}

	return plans
}

func (b *Broker) instanceExists(instanceID string) bool {
	for _, instanceCreator := range b.InstanceCreators {
		instanceExists, _ := instanceCreator.InstanceExists(instanceID)
		if instanceExists {
			return true
		}
	}
	return false
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (b *Broker) LastOperation(ctx context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}

func (b *Broker) Update(ctx context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, nil
}
