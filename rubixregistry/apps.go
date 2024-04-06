package rubixregistry

import (
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-utils-go/nversion"
	"path"
)

func (inst *RubixRegistry) GetInstalledAppVersion(appName string) string {
	rosVersion := "v1.0.0"

	files, _ := fileutils.ListFiles(path.Join(inst.AppsInstallDir, appName))
	for _, f := range files {
		if nversion.CheckVersionBool(f) {
			return f
		}
	}

	return rosVersion
}
