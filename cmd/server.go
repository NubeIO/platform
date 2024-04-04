package cmd

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-auth-go/user"
	"github.com/NubeIO/platform/config"
	"github.com/NubeIO/platform/constants"
	"github.com/NubeIO/platform/logger"
	"github.com/NubeIO/platform/router"
	"github.com/spf13/cobra"
	"os"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starting rubix-platform",
	Long:  "it starts a server for rubix-platform",
	Run:   runServer,
}

const (
	username = "admin"
	password = "admin"
)

func runServer(cmd *cobra.Command, args []string) {
	if err := config.Setup(RootCmd); err != nil {
		fmt.Errorf("error: %s", err) // here log is not setup yet...
	}
	logger.Init()
	if err := os.MkdirAll(config.Config.GetAbsDataDir(), os.FileMode(constants.Permission)); err != nil {
		panic(err)
	}
	logger.Logger.Infoln("starting edge-platform...")

	createUserIfDoesNotExist()
	r := router.Setup()

	host := "0.0.0.0"
	port := config.Config.GetPort()
	logger.Logger.Infof("server is starting at %s:%s", host, port)
	logger.Logger.Fatalf("%v", r.Run(fmt.Sprintf("%s:%s", host, port)))
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func createUserIfDoesNotExist() {
	user_, _ := user.GetUser()
	if user_ == nil {
		_, _ = user.CreateUser(&user.User{Username: username, Password: password})
	}
}
