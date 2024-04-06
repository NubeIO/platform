package appstore

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-utils-go/nversion"
	"github.com/NubeIO/platform/dto"
	"github.com/NubeIO/platform/global"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (inst *Store) GetModulesStoreModules() ([]dto.Module, error) {
	modules := make([]dto.Module, 0)

	var files []string
	err := filepath.WalkDir(global.Installer.GetModulesStorePath(), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		files = append(files, p)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileParts := strings.Split(file, "/")
		filePartsLen := len(fileParts)
		if filePartsLen > 3 {
			n := fileParts[filePartsLen-3]
			v := fileParts[filePartsLen-2]
			f := fileParts[filePartsLen-1]
			if nversion.CheckVersionBool(v) {
				arch := inst.findArch(f)
				modules = append(modules, dto.Module{Name: n, Version: v, Arch: arch})
			}
		}
	}
	return modules, err
}

func (inst *Store) UploadModuleStoreModule(app *dto.Upload) (*UploadResponse, error) {
	var file = app.File
	uploadResponse := &UploadResponse{}
	resp, err := global.Installer.Upload(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload module: %s", err.Error()))
	}
	defer os.RemoveAll(resp.TmpFile)
	uploadResponse.TmpFile = resp.TmpFile
	source := resp.UploadedFile

	fileParts := strings.Split(resp.FileName, "___")
	if len(fileParts) != 4 {
		return nil, errors.New(fmt.Sprintf("wrong module file name '%s' is being uploaded", resp.FileName))
	}
	appName := fileParts[0]
	appVersion := fileParts[1]
	fileName := fileParts[3]
	moduleStorePath := global.Installer.GetModulesStoreWithModuleVersionFolder(appName, appVersion)
	_ = os.MkdirAll(moduleStorePath, os.FileMode(global.Installer.FileMode))

	destination := path.Join(moduleStorePath, fileName)
	check := fileutils.FileExists(source)
	if !check {
		return nil, errors.New(fmt.Sprintf("upload file tmp dir not found: %s", source))
	}
	uploadResponse.UploadedFile = destination
	err = os.Rename(source, destination)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("move module error: %s", err.Error()))
	}
	uploadResponse.UploadedOk = true
	return uploadResponse, nil
}

func (inst *Store) findArch(file string) string {
	if strings.Contains(file, "armv7") {
		return "armv7"
	}
	return "amd64"
}
