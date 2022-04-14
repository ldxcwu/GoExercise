package cmd

import (
	"fmt"
	"secret"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all secrets have been managemented.",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.NewVault(encodingKey, secretsPath())
		keys, err := v.List()
		if err != nil {
			panic(err)
		}
		for _, k := range keys {
			fmt.Println(k)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
