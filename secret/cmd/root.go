package cmd

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

func init() {
	//every single time when setting or getting, with -k xxx then use the key you specificed, otherwise use the default.
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "ldxcwu", "the key to use when encoding and decoding secrets")
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets.txt")
}
