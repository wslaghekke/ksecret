package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wslaghekke/ksecret/client"
)

// editfileCmd represents the editfile command
var editfileCmd = &cobra.Command{
	Use:   "editfile [secret name] [secret key]",
	Short: "Edit a key from specified secret as file",
	Args: cobra.RangeArgs(1,2),
	Run: func(cmd *cobra.Command, args []string) {
		clientSet, clientNamespace, err := client.CreateClient(kubeconfig, namespace)
		if err != nil {
			panic(err.Error())
		}
		if len(args) == 1 {
			client.ListSecretKeys(clientSet, clientNamespace, args[0])
		} else {
			client.EditSecretKey(clientSet, clientNamespace, args[0], args[1])
		}
	},
}

func init() {
	RootCmd.AddCommand(editfileCmd)
}
