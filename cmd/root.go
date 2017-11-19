package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/wslaghekke/ksecret/client"
)

var kubeconfig string
var namespace string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ksecret",
	Short: "Kubernetes secret editor",
	Long: `Secret editor that manages base64 decoding and file editing`,
	Run: func(cmd *cobra.Command, args []string) {
		clientSet, clientNamespace, err := client.CreateClient(kubeconfig, namespace)
		if err != nil {
			panic(err.Error())
		}
		client.ListSecrets(clientSet, clientNamespace)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "config file (default is $HOME/.kube/config)")
	RootCmd.PersistentFlags().StringVar(&namespace, "namespace", "", "namespace to searchkubec")
}