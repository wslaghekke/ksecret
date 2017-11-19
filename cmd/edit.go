package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wslaghekke/ksecret/client"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
