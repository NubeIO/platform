package installer

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type AppsInstalled struct {
	ROSVersion     string         `json:"rosVersion"`
	InstalledCount int            `json:"installedCount"`
	InstalledApps  []BuildDetails `json:"installedApps"`
}

func (inst *Installer) ListAppsInstalled() (*AppsInstalled, error) {
	files, err := ioutil.ReadDir(inst.AppsInstallDir)
	if err != nil {
		return nil, err
	}
	var rosVersion string
	var installedApps []string
	var apps []BuildDetails
	for _, file := range files {
		installedApps = append(installedApps, file.Name())
	}
	for _, app := range installedApps {
		files, err = ioutil.ReadDir(fmt.Sprintf("%s/%s", inst.AppsInstallDir, app))
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			appName := app
			appVersion := "version not found"
			fileName := file.Name()
			if len(fileName) > 4 {
				if strings.Contains(fileName, "v") {
					if fileName[0:1] == "v" {
						appVersion = fileName
					}
				}
			}
			if appName == "rubix-os" {
				rosVersion = appVersion
			}
			newApp := BuildDetails{
				Name:    appName,
				Version: appVersion,
			}
			apps = append(apps, newApp)
		}
	}
	out := &AppsInstalled{
		ROSVersion:     rosVersion,
		InstalledCount: len(apps),
		InstalledApps:  apps,
	}
	return out, err
}
