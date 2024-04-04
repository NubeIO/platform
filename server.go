package platform

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type InstanceManagerHandler struct {
	*InstanceManager
	Router *gin.Engine
}

func NewInstanceManagerHandler(im *InstanceManager, router *gin.Engine) *InstanceManagerHandler {
	handler := &InstanceManagerHandler{
		InstanceManager: im,
		Router:          router,
	}

	handler.Router.POST("/api/create", handler.CreateInstance)
	handler.Router.GET("/api/start", handler.StartInstanceHandler)
	handler.Router.GET("/api/stop", handler.StopInstanceHandler)
	handler.Router.GET("/api/restart", handler.RestartInstanceHandler)
	handler.Router.GET("/api/delete", handler.DeleteInstanceHandler)
	handler.Router.GET("/api/status", handler.GetInstanceStatusHandler)
	handler.Router.GET("/api/all", handler.GetAllInstancesHandler)
	handler.Router.GET("/api/read", handler.ReadYAMLFile)
	handler.Router.GET("/api/get/pid", handler.GetPIDByPortHandler)

	handler.Router.Any("/proxy/*path", handler.ProxyHandler)

	handler.Router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))
	return handler
}

func (h *InstanceManagerHandler) CreateInstance(c *gin.Context) {
	var instance Instance
	if err := c.BindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.AddInstance(instance.Name, instance.Port)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.SaveToFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance created successfully"})
}

func (h *InstanceManagerHandler) StartInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := h.StartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance started successfully"})
}

func (h *InstanceManagerHandler) StopInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := h.StopInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance stopped successfully"})
}

func (h *InstanceManagerHandler) RestartInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := h.RestartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance restarted successfully"})
}

func (h *InstanceManagerHandler) DeleteInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := h.DeleteInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.SaveToFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance deleted successfully"})
}

func (h *InstanceManagerHandler) GetInstanceStatusHandler(c *gin.Context) {
	name := c.Query("name")
	status := h.GetInstanceStatus(name)
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *InstanceManagerHandler) GetAllInstancesHandler(c *gin.Context) {
	instances := h.GetAllInstances()
	c.JSON(http.StatusOK, gin.H{"instances": instances})
}

func (h *InstanceManagerHandler) ReadYAMLFile(c *gin.Context) {
	yamlFile, err := os.ReadFile(DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/yaml", yamlFile)
}

func (h *InstanceManagerHandler) GetPIDByPortHandler(c *gin.Context) {
	portStr := c.Query("port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable get port"})
		return
	}
	pid, err := getPID(port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pid": pid})
}

func (h *InstanceManagerHandler) ProxyHandler(c *gin.Context) {
	name := c.GetHeader("Instance-Name")
	instance, exists := h.Instances[name]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instance not found"})
		return
	}

	target := fmt.Sprintf("%s:%d", instance.IP, instance.Port)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	})

	c.Request.URL.Path = c.Param("path") // Include the endpoint path
	proxy.ServeHTTP(c.Writer, c.Request)
}
