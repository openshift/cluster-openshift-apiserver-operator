package resourcegraph

import (
	"fmt"

	"github.com/gonum/graph/encoding/dot"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/operator/resource/resourcegraph"
)

func NewResourceChainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource-graph",
		Short: "Where do resources come from? Ask your mother.",
		Run: func(cmd *cobra.Command, args []string) {
			resources := Resources()
			g := resources.NewGraph()

			data, err := dot.Marshal(g, resourcegraph.Quote("openshift-apiserver-operator"), "", "  ", false)
			if err != nil {
				klog.Fatal(err)
			}
			fmt.Println(string(data))
		},
	}

	return cmd
}

func Resources() resourcegraph.Resources {
	ret := resourcegraph.NewResources()

	payload := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "Payload", "", "cluster")).
		Add(ret)
	installer := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "Installer", "", "cluster")).
		Add(ret)
	user := resourcegraph.NewResource(resourcegraph.NewCoordinates("", "User", "", "cluster")).
		Add(ret)

	cvo := resourcegraph.NewOperator("cluster-version").
		From(payload).
		Add(ret)
	kasOperator := resourcegraph.NewOperator("kube-apiserver").
		From(cvo).
		Add(ret)
	serviceCAOperator := resourcegraph.NewOperator("service-ca").
		From(cvo).
		Add(ret)
	imageRegistryOperator := resourcegraph.NewOperator("image-registry").
		From(cvo).
		Add(ret)

	// config.openshift.io
	imageConfig := resourcegraph.NewConfig("images").
		From(user).
		From(imageRegistryOperator).
		Add(ret)
	ingressConfig := resourcegraph.NewConfig("ingresses").
		From(user).
		From(installer).
		Add(ret)
	projectConfig := resourcegraph.NewConfig("projects").
		From(user).
		Add(ret)

	// etcd certs
	fromEtcdServingCA := resourcegraph.NewConfigMap("kube-system", "etcd-serving-ca").
		Note("Static").
		From(installer).
		Add(ret)
	fromEtcdClient := resourcegraph.NewSecret("kube-system", "etcd-client").
		Note("Static").
		From(installer).
		Add(ret)
	etcdServingCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "etcd-serving-ca").
		Note("Synchronized").
		From(fromEtcdServingCA).
		Add(ret)
	etcdClient := resourcegraph.NewSecret(operatorclient.TargetNamespace, "etcd-client").
		Note("Synchronized").
		From(fromEtcdClient).
		Add(ret)

	// aggregator CA
	kasAggregatorCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "kube-apiserver-aggregator-client-ca").
		Note("Synchronized").
		From(kasOperator).
		Add(ret)
	aggregatorCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "aggregator-client-ca").
		Note("Synchronized").
		From(kasAggregatorCA).
		Add(ret)

	// client CA
	kasClientCA := resourcegraph.NewConfigMap(operatorclient.GlobalUserSpecifiedConfigNamespace, "kube-apiserver-client-ca").
		Note("Synchronized").
		From(kasOperator).
		Add(ret)
	clientCA := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "client-ca").
		Note("Synchronized").
		From(kasClientCA).
		Add(ret)

	// serving cert
	serviceCAController := resourcegraph.NewResource(resourcegraph.NewCoordinates("apps", "deployments", "openshift-service-cert-signer", "service-serving-cert-signer")).
		From(serviceCAOperator).
		Add(ret)
	servingCert := resourcegraph.NewConfigMap(operatorclient.TargetNamespace, "serving-cert").
		Note("Rotated").
		From(serviceCAController).
		Add(ret)

	// observedConfig
	config := resourcegraph.NewConfigMap(operatorclient.OperatorNamespace, "config").
		Note("Managed").
		From(imageConfig).   // to specify internal and external registries and trust
		From(ingressConfig). // routing subdomain
		From(projectConfig). // project request configuration
		Add(ret)

	// and finally our target pod
	_ = resourcegraph.NewResource(resourcegraph.NewCoordinates("", "pods", operatorclient.TargetNamespace, "openshift-apiserver")).
		From(aggregatorCA).
		From(clientCA).
		From(etcdServingCA).
		From(etcdClient).
		From(servingCert).
		From(config).
		Add(ret)

	return ret
}
