// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1

// KubeSchedulerStatusApplyConfiguration represents an declarative configuration of the KubeSchedulerStatus type for use
// with apply.
type KubeSchedulerStatusApplyConfiguration struct {
	StaticPodOperatorStatusApplyConfiguration `json:",inline"`
}

// KubeSchedulerStatusApplyConfiguration constructs an declarative configuration of the KubeSchedulerStatus type for use with
// apply.
func KubeSchedulerStatus() *KubeSchedulerStatusApplyConfiguration {
	return &KubeSchedulerStatusApplyConfiguration{}
}

// WithObservedGeneration sets the ObservedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedGeneration field is set to the value of the last call.
func (b *KubeSchedulerStatusApplyConfiguration) WithObservedGeneration(value int64) *KubeSchedulerStatusApplyConfiguration {
	b.ObservedGeneration = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *KubeSchedulerStatusApplyConfiguration) WithConditions(values ...*OperatorConditionApplyConfiguration) *KubeSchedulerStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *KubeSchedulerStatusApplyConfiguration) WithVersion(value string) *KubeSchedulerStatusApplyConfiguration {
	b.Version = &value
	return b
}

// WithReadyReplicas sets the ReadyReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadyReplicas field is set to the value of the last call.
func (b *KubeSchedulerStatusApplyConfiguration) WithReadyReplicas(value int32) *KubeSchedulerStatusApplyConfiguration {
	b.ReadyReplicas = &value
	return b
}

// WithLatestAvailableRevision sets the LatestAvailableRevision field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LatestAvailableRevision field is set to the value of the last call.
func (b *KubeSchedulerStatusApplyConfiguration) WithLatestAvailableRevision(value int32) *KubeSchedulerStatusApplyConfiguration {
	b.LatestAvailableRevision = &value
	return b
}

// WithGenerations adds the given value to the Generations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Generations field.
func (b *KubeSchedulerStatusApplyConfiguration) WithGenerations(values ...*GenerationStatusApplyConfiguration) *KubeSchedulerStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithGenerations")
		}
		b.Generations = append(b.Generations, *values[i])
	}
	return b
}

// WithLatestAvailableRevisionReason sets the LatestAvailableRevisionReason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LatestAvailableRevisionReason field is set to the value of the last call.
func (b *KubeSchedulerStatusApplyConfiguration) WithLatestAvailableRevisionReason(value string) *KubeSchedulerStatusApplyConfiguration {
	b.LatestAvailableRevisionReason = &value
	return b
}

// WithNodeStatuses adds the given value to the NodeStatuses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the NodeStatuses field.
func (b *KubeSchedulerStatusApplyConfiguration) WithNodeStatuses(values ...*NodeStatusApplyConfiguration) *KubeSchedulerStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithNodeStatuses")
		}
		b.NodeStatuses = append(b.NodeStatuses, *values[i])
	}
	return b
}
