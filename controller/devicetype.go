package controller

import (
	"github.com/NubeIO/platform/config"
	"github.com/NubeIO/platform/model"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) GetArch(c *gin.Context) {
	arch := model.Arch{Arch: config.Config.GetArch()}
	responseHandler(arch, nil, c)
}
