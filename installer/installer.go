package installer

import (
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/platform/services/rubixregistry"
	"path"
)

type Installer struct {
	RootDir         string // /data or ./
	StoreDir        string // <root_dir>/store
	TmpDir          string // <root_dir>/tmp
	BackupDir       string // <root_dir>/backup
	FileMode        int    // 0755
	DefaultTimeout  int    // 30
	AppsDownloadDir string // <root_dir>/installer/apps/download
	AppsInstallDir  string // <root_dir>/installer/apps/install
	SystemCtl       *systemctl.SystemCtl
}

func New(app *Installer, registry *rubixregistry.RubixRegistry) *Installer {
	if app == nil {
		app = &Installer{}
	}
	app.RootDir = registry.RootDir
	app.AppsDownloadDir = registry.AppsDownloadDir
	app.AppsInstallDir = registry.AppsInstallDir
	if app.FileMode == 0 {
		app.FileMode = registry.FileMode
	}
	if app.DefaultTimeout == 0 {
		app.DefaultTimeout = 30
	}
	if app.StoreDir == "" {
		app.StoreDir = path.Join(app.RootDir, "store")
	}
	if app.TmpDir == "" {
		app.TmpDir = path.Join(app.RootDir, "tmp")
	}
	if app.BackupDir == "" {
		app.BackupDir = path.Join(app.RootDir, "backup")
	}
	app.SystemCtl = systemctl.New(false, app.DefaultTimeout)
	return app
}
