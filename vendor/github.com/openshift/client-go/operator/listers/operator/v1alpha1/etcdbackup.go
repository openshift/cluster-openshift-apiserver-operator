// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/openshift/api/operator/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EtcdBackupLister helps list EtcdBackups.
// All objects returned here must be treated as read-only.
type EtcdBackupLister interface {
	// List lists all EtcdBackups in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.EtcdBackup, err error)
	// Get retrieves the EtcdBackup from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.EtcdBackup, error)
	EtcdBackupListerExpansion
}

// etcdBackupLister implements the EtcdBackupLister interface.
type etcdBackupLister struct {
	indexer cache.Indexer
}

// NewEtcdBackupLister returns a new EtcdBackupLister.
func NewEtcdBackupLister(indexer cache.Indexer) EtcdBackupLister {
	return &etcdBackupLister{indexer: indexer}
}

// List lists all EtcdBackups in the indexer.
func (s *etcdBackupLister) List(selector labels.Selector) (ret []*v1alpha1.EtcdBackup, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.EtcdBackup))
	})
	return ret, err
}

// Get retrieves the EtcdBackup from the index for a given name.
func (s *etcdBackupLister) Get(name string) (*v1alpha1.EtcdBackup, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("etcdbackup"), name)
	}
	return obj.(*v1alpha1.EtcdBackup), nil
}
