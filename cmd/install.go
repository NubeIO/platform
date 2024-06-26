package cmd

import (
	"embed"
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/platform/config"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
	"syscall"
)

var SystemdFs embed.FS

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install rubix-platform",
	Long:  "it installs rubix-platform on your device using systemd",
	Run:   install,
}

func install(cmd *cobra.Command, args []string) {
	const ServiceDir = "/lib/systemd/system"
	const ServiceDirSoftLink = "/etc/systemd/system/multi-user.target.wants"
	const ServiceFileName = "nubeio-rubix-platform.service"

	serviceFile := path.Join(ServiceDir, ServiceFileName)
	symlinkServiceFile := path.Join(ServiceDirSoftLink, ServiceFileName)

	if err := config.Setup(RootCmd); err != nil {
		fmt.Errorf("error: %s", err) // here log is not setup yet...
	}

	fmt.Println("installing edge-bios...")
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	content, err := SystemdFs.ReadFile("systemd/nubeio-rubix-platform.service")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	systemd := string(content)
	systemd = strings.Replace(systemd, "<working_dir>", wd, -1)
	fmt.Println(fmt.Sprintf("systemd file with working directory: %s", wd))

	arch := RootCmd.PersistentFlags().Lookup("arch").Value.String()
	systemd = strings.Replace(systemd, "<arch>", arch, -1)
	fmt.Println(fmt.Sprintf("systemd file with arch: %s", arch))

	fmt.Println(fmt.Sprintf("creating service file: %s...", serviceFile))
	err = os.WriteFile(serviceFile, []byte(systemd), 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("soft un-linking linux service...")
	err = syscall.Unlink(symlinkServiceFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("soft linking linux service...")
	err = syscall.Symlink(serviceFile, symlinkServiceFile)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("executing daemon-reload...")
	systemCtl := systemctl.New(false, 30)
	err = systemCtl.DaemonReload()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("enabling linux service...")
	err = systemCtl.Enable(ServiceFileName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("starting linux service...")
	err = systemCtl.Restart(ServiceFileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("successfully executed install command...")
}

func init() {
	RootCmd.AddCommand(installCmd)
}
