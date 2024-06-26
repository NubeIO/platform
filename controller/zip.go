package controller

import (
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/platform/model"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func (inst *Controller) Unzip(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")
	pathToZip := source
	if source == "" {
		responseHandler(nil, errors.New("zip source can not be empty, try /data/zip.zip"), c)
		return
	}
	if destination == "" {
		responseHandler(nil, errors.New("zip destination can not be empty, try /data/unzip-test"), c)
		return
	}
	zip, err := fileutils.Unzip(pathToZip, destination, os.FileMode(inst.FileMode))
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(zip, err, c)
}

func (inst *Controller) ZipDir(c *gin.Context) {
	source := c.Query("source")
	destination := c.Query("destination")
	pathToZip := source
	if source == "" {
		responseHandler(nil, errors.New("zip source can not be empty, try /data/flow-framework"), c)
		return
	}
	if destination == "" {
		responseHandler(nil, errors.New("zip destination can not be empty, try /data/test/flow-framework.zip"), c)
		return
	}
	exists := fileutils.DirExists(pathToZip)
	if !exists {
		responseHandler(nil, errors.New("zip source is not found"), c)
		return
	}
	err := os.MkdirAll(filepath.Dir(destination), os.FileMode(inst.FileMode))
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	err = fileutils.RecursiveZip(pathToZip, destination)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(model.Message{Message: fmt.Sprintf("zip file is created on: %s", destination)}, nil, c)
}
