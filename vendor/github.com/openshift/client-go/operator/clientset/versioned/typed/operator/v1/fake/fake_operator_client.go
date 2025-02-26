// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeOperatorV1 struct {
	*testing.Fake
}

func (c *FakeOperatorV1) Authentications() v1.AuthenticationInterface {
	return newFakeAuthentications(c)
}

func (c *FakeOperatorV1) CSISnapshotControllers() v1.CSISnapshotControllerInterface {
	return newFakeCSISnapshotControllers(c)
}

func (c *FakeOperatorV1) CloudCredentials() v1.CloudCredentialInterface {
	return newFakeCloudCredentials(c)
}

func (c *FakeOperatorV1) ClusterCSIDrivers() v1.ClusterCSIDriverInterface {
	return newFakeClusterCSIDrivers(c)
}

func (c *FakeOperatorV1) Configs() v1.ConfigInterface {
	return newFakeConfigs(c)
}

func (c *FakeOperatorV1) Consoles() v1.ConsoleInterface {
	return newFakeConsoles(c)
}

func (c *FakeOperatorV1) DNSes() v1.DNSInterface {
	return newFakeDNSes(c)
}

func (c *FakeOperatorV1) Etcds() v1.EtcdInterface {
	return newFakeEtcds(c)
}

func (c *FakeOperatorV1) IngressControllers(namespace string) v1.IngressControllerInterface {
	return newFakeIngressControllers(c, namespace)
}

func (c *FakeOperatorV1) InsightsOperators() v1.InsightsOperatorInterface {
	return newFakeInsightsOperators(c)
}

func (c *FakeOperatorV1) KubeAPIServers() v1.KubeAPIServerInterface {
	return newFakeKubeAPIServers(c)
}

func (c *FakeOperatorV1) KubeControllerManagers() v1.KubeControllerManagerInterface {
	return newFakeKubeControllerManagers(c)
}

func (c *FakeOperatorV1) KubeSchedulers() v1.KubeSchedulerInterface {
	return newFakeKubeSchedulers(c)
}

func (c *FakeOperatorV1) KubeStorageVersionMigrators() v1.KubeStorageVersionMigratorInterface {
	return newFakeKubeStorageVersionMigrators(c)
}

func (c *FakeOperatorV1) MachineConfigurations() v1.MachineConfigurationInterface {
	return newFakeMachineConfigurations(c)
}

func (c *FakeOperatorV1) Networks() v1.NetworkInterface {
	return newFakeNetworks(c)
}

func (c *FakeOperatorV1) OLMs() v1.OLMInterface {
	return newFakeOLMs(c)
}

func (c *FakeOperatorV1) OpenShiftAPIServers() v1.OpenShiftAPIServerInterface {
	return newFakeOpenShiftAPIServers(c)
}

func (c *FakeOperatorV1) OpenShiftControllerManagers() v1.OpenShiftControllerManagerInterface {
	return newFakeOpenShiftControllerManagers(c)
}

func (c *FakeOperatorV1) ServiceCAs() v1.ServiceCAInterface {
	return newFakeServiceCAs(c)
}

func (c *FakeOperatorV1) ServiceCatalogAPIServers() v1.ServiceCatalogAPIServerInterface {
	return newFakeServiceCatalogAPIServers(c)
}

func (c *FakeOperatorV1) ServiceCatalogControllerManagers() v1.ServiceCatalogControllerManagerInterface {
	return newFakeServiceCatalogControllerManagers(c)
}

func (c *FakeOperatorV1) Storages() v1.StorageInterface {
	return newFakeStorages(c)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeOperatorV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
