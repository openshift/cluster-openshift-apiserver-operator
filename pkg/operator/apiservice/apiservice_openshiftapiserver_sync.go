package apiservice

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog/v2"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	apiregistrationv1lister "k8s.io/kube-aggregator/pkg/client/listers/apiregistration/v1"

	operatorlistersv1 "github.com/openshift/client-go/operator/listers/operator/v1"
	"github.com/openshift/library-go/pkg/operator/events"
)

// APIServicesToMange preserve state and clients required to return an authoritative list of API services this operate must manage
type APIServicesToManage struct {
	authOperatorLister                         operatorlistersv1.AuthenticationLister
	apiregistrationv1Lister                    apiregistrationv1lister.APIServiceLister
	allPossibleAPIServices                     []*apiregistrationv1.APIService
	eventRecorder                              events.Recorder
	apiGroupsManagedByExternalServer           sets.String
	apiGroupsManagedByExternalServerAnnotation string
	currentAPIServicesToManage                 []*apiregistrationv1.APIService
}

// NewAPIServicesToManage returns an object that knows how to construct an authoritative list of API services this operate must manage
func NewAPIServicesToManage(
	apiregistrationv1Lister apiregistrationv1lister.APIServiceLister,
	authOperatorLister operatorlistersv1.AuthenticationLister,
	allPossibleAPIServices []*apiregistrationv1.APIService,
	eventRecorder events.Recorder,
	apiGroupsManagedByExternalServer sets.String,
	apiGroupsManagedByExternalServerAnnotation string) *APIServicesToManage {
	return &APIServicesToManage{
		authOperatorLister:                         authOperatorLister,
		apiregistrationv1Lister:                    apiregistrationv1Lister,
		allPossibleAPIServices:                     allPossibleAPIServices,
		eventRecorder:                              eventRecorder,
		apiGroupsManagedByExternalServer:           apiGroupsManagedByExternalServer,
		apiGroupsManagedByExternalServerAnnotation: apiGroupsManagedByExternalServerAnnotation,
		currentAPIServicesToManage:                 allPossibleAPIServices,
	}
}

// GetAPIServicesToManage returns the desired list of API Services that will be managed by this operator
// note that some services might be managed by an external operators/servers
func (a *APIServicesToManage) GetAPIServicesToManage() ([]*apiregistrationv1.APIService, error) {
	if externalOperatorPreconditionErr := a.externalOperatorPrecondition(); externalOperatorPreconditionErr != nil {
		klog.V(4).Infof("unable to determine if an external operator should take OAuth APIs over due to %v, returning authoritative/initial API Services list", externalOperatorPreconditionErr)
		return a.allPossibleAPIServices, nil
	}

	newAPIServicesToManage := []*apiregistrationv1.APIService{}
	for _, apiService := range a.allPossibleAPIServices {
		if a.apiGroupsManagedByExternalServer.Has(apiService.Name) && a.isAPIServiceAnnotatedByExternalServer(apiService) {
			continue
		}
		newAPIServicesToManage = append(newAPIServicesToManage, apiService)
	}

	if changed, newAPIServicesSet := apiServicesChanged(a.currentAPIServicesToManage, newAPIServicesToManage); changed {
		a.eventRecorder.Eventf("APIServicesToManageChanged", "The new API Services list this operator will manage is %v", newAPIServicesSet.List())
	}

	a.currentAPIServicesToManage = newAPIServicesToManage
	return a.currentAPIServicesToManage, nil
}

func (a *APIServicesToManage) isAPIServiceAnnotatedByExternalServer(apiService *apiregistrationv1.APIService) bool {
	existingApiService, err := a.apiregistrationv1Lister.Get(apiService.Name)
	if err != nil {
		a.eventRecorder.Warningf("APIServicesToManageAnnotation", "unable to determine if the following API Service %s was annotated by an external operator (it should be) due to %v", apiService.Name, err)
		return false
	}

	if _, ok := existingApiService.Annotations[a.apiGroupsManagedByExternalServerAnnotation]; ok {
		return true

	}
	return false
}

// externalOperatorPrecondition checks whether authentication operator will manage OAuth API Resources by checking ManagingOAuthAPIServer status field
func (a *APIServicesToManage) externalOperatorPrecondition() error {
	authOperator, err := a.authOperatorLister.Get("cluster")
	if err != nil {
		return err
	}

	if !authOperator.Status.ManagingOAuthAPIServer {
		return fmt.Errorf("%q status field set to false", "ManagingOAuthAPIServer")
	}

	return nil
}

func apiServicesChanged(old []*apiregistrationv1.APIService, new []*apiregistrationv1.APIService) (bool, sets.String) {
	oldSet := sets.String{}
	for _, oldService := range old {
		oldSet.Insert(oldService.Name)
	}

	newSet := sets.String{}
	for _, newService := range new {
		newSet.Insert(newService.Name)
	}

	removed := oldSet.Difference(newSet).List()
	added := newSet.Difference(oldSet).List()
	return len(removed) > 0 || len(added) > 0, newSet
}
