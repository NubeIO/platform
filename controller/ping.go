package controller

import (
	"github.com/NubeIO/platform/interfaces"
	"github.com/NubeIO/platform/utils/release"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) Ping(c *gin.Context) {
	responseHandler(interfaces.Ping{Version: release.GetVersion()}, nil, c)
}
