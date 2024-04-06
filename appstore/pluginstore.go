package appstore

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/platform/dto"
	"github.com/NubeIO/platform/global"
	"github.com/NubeIO/platform/installer"
	"io/ioutil"
	"os"
)

func (inst *Store) GetPluginsStorePlugins() ([]installer.BuildDetails, error) {
	pluginStore := global.Installer.GetPluginsStorePath()
	_ = os.MkdirAll(pluginStore, os.FileMode(global.Installer.FileMode))
	files, err := ioutil.ReadDir(pluginStore)
	if err != nil {
		return nil, err
	}
	plugins := make([]installer.BuildDetails, 0)
	for _, file := range files {
		plugins = append(plugins, *global.Installer.GetZipBuildDetails(file.Name()))
	}
	return plugins, err
}

func (inst *Store) UploadPluginStorePlugin(app *dto.Upload) (*UploadResponse, error) {
	var file = app.File
	uploadResponse := &UploadResponse{}
	resp, err := global.Installer.Upload(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("upload plugin: %s", err.Error()))
	}
	defer os.RemoveAll(resp.TmpFile)
	uploadResponse.TmpFile = resp.TmpFile
	source := resp.UploadedFile

	destination := global.Installer.GetPluginsStoreWithFile(resp.FileName)
	check := fileutils.FileExists(source)
	if !check {
		return nil, errors.New(fmt.Sprintf("upload file tmp dir not found: %s", source))
	}
	uploadResponse.UploadedFile = destination
	err = os.Rename(source, destination)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("move plugin error: %s", err.Error()))
	}
	uploadResponse.UploadedOk = true
	return uploadResponse, nil
}
