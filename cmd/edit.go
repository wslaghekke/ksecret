package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wslaghekke/ksecret/client"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit [secret name]",
	Short: "A brief description of your command",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clientSet, clientNamespace, err := client.CreateClient(kubeconfig, namespace)
		if err != nil {
			panic(err.Error())
		}
		client.EditSecret(clientSet, clientNamespace, args[0])
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
}
