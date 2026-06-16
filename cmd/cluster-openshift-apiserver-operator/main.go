package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	"k8s.io/component-base/cli"

	kmshealth "github.com/openshift/library-go/pkg/operator/encryption/kms/health"
	"github.com/openshift/library-go/pkg/operator/v1helpers"

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
	cmd.AddCommand(kmshealth.NewCommand(context.Background(), func(config *rest.Config) (v1helpers.OperatorClient, error) {
		// TODO: replace with a real operator client once the health reporter's condition writer
		// is implemented in library-go.
		return nil, nil
	}))

	return cmd
}
