package controller

import (
	"github.com/NubeIO/lib-dhcpd/dhcpd"
	"github.com/NubeIO/platform/services/info"
	"github.com/gin-gonic/gin"
)

func (inst *Controller) InternetIP(c *gin.Context) {
	data, err := inst.SystemInfo.GetInternetIP()
	responseHandler(data, err, c)
}

func (inst *Controller) RestartNetworking(c *gin.Context) {
	data, err := inst.Networking.RestartNetworking()
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceUpDown(c *gin.Context) {
	var m info.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.InterfaceUpDown(m)
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceUp(c *gin.Context) {
	var m info.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.InterfaceUp(m)
	responseHandler(data, err, c)
}

func (inst *Controller) InterfaceDown(c *gin.Context) {
	var m info.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.InterfaceDown(m)
	responseHandler(data, err, c)
}

func (inst *Controller) DHCPPortExists(c *gin.Context) {
	var m info.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.DHCPPortExists(m)
	responseHandler(data, err, c)
}

func (inst *Controller) DHCPSetAsAuto(c *gin.Context) {
	var m info.NetworkingBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.DHCPSetAsAuto(m)
	responseHandler(data, err, c)
}

func (inst *Controller) DHCPSetStaticIP(c *gin.Context) {
	var m *dhcpd.SetStaticIP
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.DHCPSetStaticIP(m)
	responseHandler(data, err, c)
}

func (inst *Controller) UWFActive(c *gin.Context) {
	data, err := inst.Networking.UWFActive()
	responseHandler(data, err, c)
}

func (inst *Controller) UWFEnable(c *gin.Context) {
	data, err := inst.Networking.UWFEnable()
	responseHandler(data, err, c)
}

func (inst *Controller) UWFDisable(c *gin.Context) {
	data, err := inst.Networking.UWFDisable()
	responseHandler(data, err, c)
}

func (inst *Controller) UWFStatus(c *gin.Context) {
	data, err := inst.Networking.UWFStatus()
	responseHandler(data, err, c)
}

func (inst *Controller) UWFStatusList(c *gin.Context) {
	data, err := inst.Networking.UWFStatusList()
	responseHandler(data, err, c)
}

func (inst *Controller) UWFOpenPort(c *gin.Context) {
	var m info.UFWBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.UWFOpenPort(m)
	responseHandler(data, err, c)
}

func (inst *Controller) UWFClosePort(c *gin.Context) {
	var m info.UFWBody
	err := c.ShouldBindJSON(&m)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	data, err := inst.Networking.UWFClosePort(m)
	responseHandler(data, err, c)
}
