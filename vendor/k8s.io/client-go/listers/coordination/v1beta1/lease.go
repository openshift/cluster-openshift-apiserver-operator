/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// LeaseLister helps list Leases.
// All objects returned here must be treated as read-only.
type LeaseLister interface {
	// List lists all Leases in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*coordinationv1beta1.Lease, err error)
	// Leases returns an object that can list and get Leases.
	Leases(namespace string) LeaseNamespaceLister
	LeaseListerExpansion
}

// leaseLister implements the LeaseLister interface.
type leaseLister struct {
	listers.ResourceIndexer[*coordinationv1beta1.Lease]
}

// NewLeaseLister returns a new LeaseLister.
func NewLeaseLister(indexer cache.Indexer) LeaseLister {
	return &leaseLister{listers.New[*coordinationv1beta1.Lease](indexer, coordinationv1beta1.Resource("lease"))}
}

// Leases returns an object that can list and get Leases.
func (s *leaseLister) Leases(namespace string) LeaseNamespaceLister {
	return leaseNamespaceLister{listers.NewNamespaced[*coordinationv1beta1.Lease](s.ResourceIndexer, namespace)}
}

// LeaseNamespaceLister helps list and get Leases.
// All objects returned here must be treated as read-only.
type LeaseNamespaceLister interface {
	// List lists all Leases in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*coordinationv1beta1.Lease, err error)
	// Get retrieves the Lease from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*coordinationv1beta1.Lease, error)
	LeaseNamespaceListerExpansion
}

// leaseNamespaceLister implements the LeaseNamespaceLister
// interface.
type leaseNamespaceLister struct {
	listers.ResourceIndexer[*coordinationv1beta1.Lease]
}
