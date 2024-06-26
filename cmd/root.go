package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nube-cli",
	Short: "description",
	Long:  `description`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
	}
}

var flgRoot struct {
	prod      bool
	auth      bool
	port      int
	rootDir   string
	appDir    string
	dataDir   string
	configDir string
	arch      string
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&flgRoot.prod, "prod", "", false, "prod")
	RootCmd.PersistentFlags().BoolVarP(&flgRoot.auth, "auth", "", true, "auth")
	RootCmd.PersistentFlags().IntVarP(&flgRoot.port, "port", "p", 1772, "port (default 1772)")
	RootCmd.PersistentFlags().StringVarP(&flgRoot.rootDir, "root-dir", "r", "./", "root dir") // in production it will be `/data`
	RootCmd.PersistentFlags().StringVarP(&flgRoot.appDir, "app-dir", "a", "./", "app dir")    // in production it will be `rubix-bios`
	RootCmd.PersistentFlags().StringVarP(&flgRoot.dataDir, "data-dir", "d", "data", "data dir")
	RootCmd.PersistentFlags().StringVarP(&flgRoot.configDir, "config-dir", "c", "config", "config dir")
	RootCmd.PersistentFlags().StringVarP(&flgRoot.arch, "arch", "", "armv7", "device type")
}
