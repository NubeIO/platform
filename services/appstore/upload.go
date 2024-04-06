package appstore

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/platform/dto"
	"os"
	"path"
)

type UploadResponse struct {
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	UploadedOk   bool   `json:"uploadedOk,omitempty"`
	TmpFile      string `json:"tmpFile,omitempty"`
	UploadedFile string `json:"uploadedFile,omitempty"`
}

func (inst *Store) UploadAddOnAppStore(app *dto.Upload) (*UploadResponse, error) {
	if app.Name == "" {
		return nil, errors.New("app_name can not be empty")
	}
	if app.Version == "" {
		return nil, errors.New("app_version can not be empty")
	}
	if app.Arch == "" {
		return nil, errors.New("arch_type can not be empty, try armv7 amd64")
	}
	err := os.MkdirAll(inst.Installer.GetAppsStoreAppPathWithArchVersion(app.Name, app.Arch, app.Version), os.FileMode(inst.Installer.FileMode))
	if err != nil {
		return nil, err
	}
	var file = app.File
	uploadResp := &UploadResponse{
		Name:         app.Name,
		Version:      app.Version,
		UploadedOk:   false,
		TmpFile:      "",
		UploadedFile: "",
	}
	resp, err := inst.Installer.Upload(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload app: %s", err.Error()))
	}
	defer os.RemoveAll(resp.TmpFile)
	uploadResp.TmpFile = resp.TmpFile
	source := resp.UploadedFile
	destination := path.Join(inst.Installer.GetAppsStoreAppPathWithArchVersion(app.Name, app.Arch, app.Version), resp.FileName)
	check := fileutils.FileExists(source)
	if !check {
		return nil, errors.New(fmt.Sprintf("upload file tmp dir not found: %s", source))
	}
	uploadResp.UploadedFile = destination
	err = os.Rename(source, destination)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("move build error: %s", err.Error()))
	}
	uploadResp.UploadedOk = true
	return uploadResp, nil
}
