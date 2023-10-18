// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1 "github.com/openshift/api/operator/v1"
	operatorv1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEtcds implements EtcdInterface
type FakeEtcds struct {
	Fake *FakeOperatorV1
}

var etcdsResource = v1.SchemeGroupVersion.WithResource("etcds")

var etcdsKind = v1.SchemeGroupVersion.WithKind("Etcd")

// Get takes name of the etcd, and returns the corresponding etcd object, and an error if there is any.
func (c *FakeEtcds) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Etcd, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(etcdsResource, name), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// List takes label and field selectors, and returns the list of Etcds that match those selectors.
func (c *FakeEtcds) List(ctx context.Context, opts metav1.ListOptions) (result *v1.EtcdList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(etcdsResource, etcdsKind, opts), &v1.EtcdList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.EtcdList{ListMeta: obj.(*v1.EtcdList).ListMeta}
	for _, item := range obj.(*v1.EtcdList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested etcds.
func (c *FakeEtcds) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(etcdsResource, opts))
}

// Create takes the representation of a etcd and creates it.  Returns the server's representation of the etcd, and an error, if there is any.
func (c *FakeEtcds) Create(ctx context.Context, etcd *v1.Etcd, opts metav1.CreateOptions) (result *v1.Etcd, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(etcdsResource, etcd), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// Update takes the representation of a etcd and updates it. Returns the server's representation of the etcd, and an error, if there is any.
func (c *FakeEtcds) Update(ctx context.Context, etcd *v1.Etcd, opts metav1.UpdateOptions) (result *v1.Etcd, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(etcdsResource, etcd), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEtcds) UpdateStatus(ctx context.Context, etcd *v1.Etcd, opts metav1.UpdateOptions) (*v1.Etcd, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(etcdsResource, "status", etcd), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// Delete takes name of the etcd and deletes it. Returns an error if one occurs.
func (c *FakeEtcds) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(etcdsResource, name, opts), &v1.Etcd{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEtcds) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(etcdsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1.EtcdList{})
	return err
}

// Patch applies the patch and returns the patched etcd.
func (c *FakeEtcds) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Etcd, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdsResource, name, pt, data, subresources...), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied etcd.
func (c *FakeEtcds) Apply(ctx context.Context, etcd *operatorv1.EtcdApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Etcd, err error) {
	if etcd == nil {
		return nil, fmt.Errorf("etcd provided to Apply must not be nil")
	}
	data, err := json.Marshal(etcd)
	if err != nil {
		return nil, err
	}
	name := etcd.Name
	if name == nil {
		return nil, fmt.Errorf("etcd.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdsResource, *name, types.ApplyPatchType, data), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeEtcds) ApplyStatus(ctx context.Context, etcd *operatorv1.EtcdApplyConfiguration, opts metav1.ApplyOptions) (result *v1.Etcd, err error) {
	if etcd == nil {
		return nil, fmt.Errorf("etcd provided to Apply must not be nil")
	}
	data, err := json.Marshal(etcd)
	if err != nil {
		return nil, err
	}
	name := etcd.Name
	if name == nil {
		return nil, fmt.Errorf("etcd.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdsResource, *name, types.ApplyPatchType, data, "status"), &v1.Etcd{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.Etcd), err
}
