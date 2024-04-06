package controller

import (
	"fmt"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/platform/config"
	"github.com/NubeIO/platform/model"
	"github.com/NubeIO/platform/services/info"
	systeminfo "github.com/NubeIO/platform/services/system"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Controller struct {
	SystemCtl  *systemctl.SystemCtl
	FileMode   int
	Instances  map[string]*Instance
	Lock       sync.Mutex
	Config     *config.Configuration
	SystemInfo systeminfo.System
	Networking *info.System
}

type Response struct {
	StatusCode   int         `json:"status_code"`
	ErrorMessage string      `json:"error_message"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}

func responseHandler(body interface{}, err error, c *gin.Context, statusCode ...int) {
	var code int
	if err != nil {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusBadRequest
		}
		msg := model.Message{
			Message: fmt.Sprintf("platform: %s", err.Error()),
		}
		c.JSON(code, msg)
	} else {
		if len(statusCode) > 0 {
			code = statusCode[0]
		} else {
			code = http.StatusOK
			if c.Request.Method == "POST" {
				code = http.StatusCreated
			}
		}
		c.JSON(code, body)
	}
}

func Builder(ip string, port int) (*url.URL, error) {
	return url.ParseRequestURI(CheckHTTP(fmt.Sprintf("%s:%d", ip, port)))
}

func CheckHTTP(address string) string {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		return "http://" + address
	}
	return address
}
