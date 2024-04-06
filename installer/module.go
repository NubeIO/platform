package installer

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-utils-go/nstring"
	"path"
	"strings"
)

func (inst *Installer) GetModulesStoreModule(module, version, arch string) (*string, error) {
	moduleLocation := inst.GetModulesStoreWithModuleVersionFolder(module, version)
	files, err := fileutils.ListFiles(moduleLocation)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.Contains(file, arch) {
			return nstring.New(path.Join(moduleLocation, file)), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no files found in %s for arch %s", moduleLocation, arch))
}
