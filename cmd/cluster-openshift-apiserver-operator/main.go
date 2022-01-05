package main

import (
	"os"

	"github.com/spf13/cobra"

	"k8s.io/component-base/cli"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/cmd/operator"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/cmd/resourcegraph"
)

func main() {
	command := NewSSCSCommand()
	code := cli.Run(command)
	os.Exit(code)
}

func NewSSCSCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster-openshift-apiserver-operator",
		Short: "OpenShift cluster openshift-apiserver operator",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	cmd.AddCommand(operator.NewOperator())
	cmd.AddCommand(resourcegraph.NewResourceChainCommand())

	return cmd
}
