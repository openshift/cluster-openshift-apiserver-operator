package main

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/component-base/cli"

	kmshealth "github.com/openshift/library-go/pkg/operator/encryption/kms/health"
	kmspreflight "github.com/openshift/library-go/pkg/operator/encryption/kms/preflight"

	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/cmd/operator"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/cmd/resourcegraph"
	"github.com/openshift/cluster-openshift-apiserver-operator/pkg/operator/encryptionstatusprovider"
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
	cmd.AddCommand(kmshealth.NewCommand(context.Background(), encryptionstatusprovider.NewOpenShiftAPIServerEncryptionStatusProviderFromConfig))
	cmd.AddCommand(kmspreflight.NewCommand(context.Background()))

	return cmd
}
