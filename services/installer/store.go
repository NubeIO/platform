package installer

import (
	"path"
)

func (inst *Installer) GetAppsStorePath() string {
	return path.Join(inst.StoreDir, "apps") // <root_dir>/store/apps
}

func (inst *Installer) GetAppsStoreAppPath(appName string) string {
	return path.Join(inst.GetAppsStorePath(), appName) // <root_dir>//store/apps/<app_name>
}

func (inst *Installer) GetAppsStoreAppPathWithArchVersion(appName, arch, version string) string {
	return path.Join(inst.GetAppsStoreAppPath(appName), arch, version) // <root_dir>//store/apps/<app_name>/<arch>/<version>
}

func (inst *Installer) GetPluginsStorePath() string {
	return path.Join(inst.StoreDir, "plugins") // <root_dir>//store/plugins
}

func (inst *Installer) GetPluginsStoreWithFile(fileName string) string {
	p := path.Join(inst.GetPluginsStorePath(), fileName) // <root_dir>//store/plugins/<plugin_file>
	return p
}

func (inst *Installer) GetModulesStorePath() string {
	return path.Join(inst.StoreDir, "modules") // <root_dir>//store/modules
}

func (inst *Installer) GetModulesStoreWithModuleVersionFolder(name, version string) string {
	return path.Join(inst.StoreDir, "modules", name, version) // <root_dir>//store/modules/<name>/<version>
}
