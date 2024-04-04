package controller

import (
	"github.com/NubeIO/platform/model"
	"github.com/gin-gonic/gin"
	"os/exec"
)

func (inst *Controller) RebootHost(c *gin.Context) {
	cmd := exec.Command("shutdown", "-r", "now")
	_, err := cmd.Output()
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(model.Message{Message: "restarted the device successfully"}, nil, c)
}
