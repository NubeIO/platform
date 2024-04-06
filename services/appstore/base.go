package appstore

import (
	"github.com/NubeIO/platform/services/installer"
	log "github.com/sirupsen/logrus"
	"os"
)

type StoreDatabase interface {
}
type Store struct {
	Installer *installer.Installer
}

func New(store *Store) *Store {
	err := store.initMakeAllDirs()
	if err != nil {
		log.Fatal(err)
	}
	store.Installer = installer.New(nil, nil)
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
