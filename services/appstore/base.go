package appstore

import (
	"github.com/NubeIO/platform/services/installer"
	"github.com/NubeIO/platform/services/rubixregistry"
	log "github.com/sirupsen/logrus"
	"os"
)

type StoreDatabase interface {
}
type Store struct {
	Installer *installer.Installer
}

func New(rootDir string) *Store {
	registry := rubixregistry.New(rootDir)
	store := &Store{
		Installer: installer.New(&installer.Installer{}, registry),
	}
	err := store.initMakeAllDirs()
	if err != nil {
		log.Fatal(err)
	}
	return store
}

func (inst *Store) initMakeAllDirs() error {
	if err := os.MkdirAll(inst.Installer.GetAppsStorePath(), os.FileMode(inst.Installer.FileMode)); err != nil {
		return err
	}
	if err := os.MkdirAll(inst.Installer.GetPluginsStorePath(), os.FileMode(inst.Installer.FileMode)); err != nil {
		return err
	}
	if err := os.MkdirAll(inst.Installer.GetModulesStorePath(), os.FileMode(inst.Installer.FileMode)); err != nil {
		return err
	}
	return nil
}
