package operator

import (
	"github.com/spf13/cobra"

	"k8s.io/utils/clock"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/version"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
)

func NewOperator() *cobra.Command {
	cmd := controllercmd.
		NewControllerCommandConfig("openshift-apiserver-operator", version.Get(), operator.RunOperator, clock.RealClock{}).
		NewCommand()
	cmd.Use = "operator"
	cmd.Short = "Start the Cluster openshift-apiserver Operator"

	return cmd
}
