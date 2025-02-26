// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/openshift/api/operator/v1"
	operatorv1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1"
	typedoperatorv1 "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	gentype "k8s.io/client-go/gentype"
)

// fakeCSISnapshotControllers implements CSISnapshotControllerInterface
type fakeCSISnapshotControllers struct {
	*gentype.FakeClientWithListAndApply[*v1.CSISnapshotController, *v1.CSISnapshotControllerList, *operatorv1.CSISnapshotControllerApplyConfiguration]
	Fake *FakeOperatorV1
}

func newFakeCSISnapshotControllers(fake *FakeOperatorV1) typedoperatorv1.CSISnapshotControllerInterface {
	return &fakeCSISnapshotControllers{
		gentype.NewFakeClientWithListAndApply[*v1.CSISnapshotController, *v1.CSISnapshotControllerList, *operatorv1.CSISnapshotControllerApplyConfiguration](
			fake.Fake,
			"",
			v1.SchemeGroupVersion.WithResource("csisnapshotcontrollers"),
			v1.SchemeGroupVersion.WithKind("CSISnapshotController"),
			func() *v1.CSISnapshotController { return &v1.CSISnapshotController{} },
			func() *v1.CSISnapshotControllerList { return &v1.CSISnapshotControllerList{} },
			func(dst, src *v1.CSISnapshotControllerList) { dst.ListMeta = src.ListMeta },
			func(list *v1.CSISnapshotControllerList) []*v1.CSISnapshotController {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1.CSISnapshotControllerList, items []*v1.CSISnapshotController) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
