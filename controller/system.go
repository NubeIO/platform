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

func (inst *Controller) GetSystemInfo(c *gin.Context) {
	queryParams := c.Request.URL.Query() //eg: /api/system/info?uptime&&ip
	var args []string
	for key, _ := range queryParams {
		args = append(args, key)
	}
	methods, err := inst.SystemInfo.ExecuteMethods(args)
	responseHandler(methods, err, c)
}
