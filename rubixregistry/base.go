package rubixregistry

import (
	"path"
)

type RubixRegistry struct {
	RootDir              string // /data
	RegistryDir          string // /<root_dir>/rubix-registry
	LegacyDeviceInfoFile string // /<root_dir>/rubix-registry/device_info.json // TODO: remove after migration done
	GlobalUUIDFile       string // /<root_dir>/rubix-registry/global_uuid.txt
	FileMode             int    // 0755
	ProductInfoPath      string // /<root_dir>/product.json
	AppsDownloadDir      string // <root_dir>/installer/apps/download
	AppsInstallDir       string // <root_dir>/installer/apps/install
}

func New(rootDir string) *RubixRegistry {
	registry := &RubixRegistry{
		RootDir:              rootDir,
		RegistryDir:          path.Join(rootDir, "rubix-registry"),
		LegacyDeviceInfoFile: path.Join(rootDir, "rubix-registry/device_info.json"), // TODO: remove after migration done
		GlobalUUIDFile:       path.Join(rootDir, "rubix-registry/global_uuid.txt"),
		FileMode:             0755,
		ProductInfoPath:      path.Join(rootDir, "product.json"),
		AppsDownloadDir:      path.Join(rootDir, "installer/apps/download"),
		AppsInstallDir:       path.Join(rootDir, "installer/apps/install"),
	}
	return registry
}
