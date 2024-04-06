package installer

import (
	"fmt"
	"github.com/NubeIO/lib-utils-go/nuuid"
	"github.com/NubeIO/platform/constants"
	"os"
	"path"
	"time"
)

func (inst *Installer) GetAppDataPath(appName string) string {
	dataDirName := constants.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName) // <root_dir>/rubix-wires
}

func (inst *Installer) GetAppDataDataPath(appName string) string {
	dataDirName := constants.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName, "data") // <root_dir>/rubix-wires/data
}

func (inst *Installer) GetAppDataConfigPath(appName string) string {
	dataDirName := constants.GetDataDirNameFromAppName(appName)
	return path.Join(inst.RootDir, dataDirName, "config") // <root_dir>/rubix-wires/config
}

func (inst *Installer) GetAppInstallPath(appName string) string {
	repoName := constants.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName) // <root_dir>/installer/apps/install/wires-builds
}

func (inst *Installer) GetAppInstallPathWithVersion(appName, version string) string {
	repoName := constants.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsInstallDir, repoName, version) // <root_dir>/installer/apps/install/wires-builds/v0.0.1
}

func (inst *Installer) GetAppDownloadPath(appName string) string {
	repoName := constants.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName) // <root_dir>/installer/apps/download/wires-builds
}

func (inst *Installer) GetAppDownloadPathWithVersion(appName, version string) string {
	repoName := constants.GetRepoNameFromAppName(appName)
	return path.Join(inst.AppsDownloadDir, repoName, version) // <root_dir>/installer/apps/download/wires-builds/v0.0.1
}

func (inst *Installer) GetEmptyNewTmpFolder() string {
	return path.Join(inst.TmpDir, nuuid.ShortUUID("tmp")) // <root_dir>/tmp/tmp_45EA34EB
}

func (inst *Installer) MakeTmpDir() error {
	return os.MkdirAll(inst.TmpDir, os.FileMode(inst.FileMode)) // <root_dir>/tmp
}

func (inst *Installer) MakeTmpDirUpload() (string, error) {
	tmpDir := inst.GetEmptyNewTmpFolder()
	err := os.MkdirAll(tmpDir, os.FileMode(inst.FileMode)) // <root_dir>/tmp/tmp_45EA34EB
	return tmpDir, err
}

func (inst *Installer) GetROSPluginDownloadPath() string {
	repoName := constants.GetRepoNameFromAppName(constants.RubixOs)
	return path.Join(inst.AppsDownloadDir, repoName, "plugins") // <root_dir>/installer/apps/download/rubix-os/plugins
}

func (inst *Installer) GetROSPluginInstallPath() string {
	return path.Join(inst.GetAppDataDataPath(constants.RubixOs), "plugins") // <root_dir>/rubix-os/data/plugins
}

func (inst *Installer) GetROSModuleDownloadPath() string {
	repoName := constants.GetRepoNameFromAppName(constants.RubixOs)
	return path.Join(inst.AppsDownloadDir, repoName, "modules") // <root_dir>/installer/apps/download/rubix-os/modules
}

func (inst *Installer) GetROSModuleDownloadWithModuleVersionPath(name, version string) string {
	return path.Join(inst.GetROSModuleDownloadPath(), name, version) // <root_dir>/installer/apps/download/rubix-os/modules/<module>/<version>
}

func (inst *Installer) GetAppBackupPath(appName, version string) string {
	return path.Join(inst.BackupDir, appName,
		fmt.Sprintf("%s_%s", time.Now().UTC().Format("20060102150405"), version)) // <root_dir>/backup/rubix-wires/<time_value>_v0.0.1
}
