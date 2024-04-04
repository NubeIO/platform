package main

import (
	"embed"
	"github.com/NubeIO/platform/cmd"
	"github.com/NubeIO/platform/release"
)

//go:embed systemd/*
var systemdFs embed.FS

//go:embed VERSION
var versionFs embed.FS

func main() {
	cmd.SystemdFs = systemdFs
	release.VersionFs = versionFs
	cmd.Execute()
}
