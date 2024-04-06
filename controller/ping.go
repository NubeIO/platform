package controller

import (
	"github.com/NubeIO/platform/model"
	"github.com/NubeIO/platform/utils/release"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) Ping(c *gin.Context) {
	responseHandler(model.Ping{Version: release.GetVersion()}, nil, c)
}
