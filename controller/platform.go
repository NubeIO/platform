package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const DB = "db.yaml"
const BinaryName = "main"

type Instance struct {
	Name  string `yaml:"name"`
	IP    string `yaml:"ip"`
	Port  int    `yaml:"port"`
	PID   int    `yaml:"pid"`
	Error string `yaml:"error"`
}

func (inst *Controller) LoadFromFile(filePath string) error {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	inst.Instances = make(map[string]*Instance)
	err = yaml.Unmarshal(yamlFile, &inst.Instances)
	if err != nil {
		return err
	}

	return nil
}

func (inst *Controller) SaveToFile() error {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()

	yamlData, err := yaml.Marshal(&inst.Instances)
	if err != nil {
		return err
	}

	err = os.WriteFile(DB, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (inst *Controller) AddInstance(name string, port int) error {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()

	_, exists := inst.Instances[name]
	if exists {
		return fmt.Errorf("instance with name %s already exists", name)
	}

	for _, instance := range inst.Instances {
		if instance.Port == port {
			return fmt.Errorf("instance with port %d already exists", port)
		}
	}

	inst.Instances[name] = &Instance{Name: name, Port: port}
	err := inst.StartInstance(name)
	if err != nil {
		return err
	}
	return nil
}

func (inst *Controller) DeleteInstance(name string) error {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()

	_, exists := inst.Instances[name]
	if !exists {
		return fmt.Errorf("instance with name %s not found", name)
	}

	delete(inst.Instances, name)
	err := inst.StopInstance(name)
	if err != nil {
		return err
	}
	return nil
}

func (inst *Controller) StartInstance(name string) error {
	cmd := exec.Command(fmt.Sprintf("./%s", BinaryName), fmt.Sprintf("-p=%d", inst.Instances[name].Port))
	err := cmd.Start()
	if err != nil {
		return err
	}
	inst.Instances[name].PID = cmd.Process.Pid
	return nil
}

func (inst *Controller) StopInstance(name string) error {
	port := inst.Instances[name].Port
	pid, err := getPID(port)
	if err != nil {
		return fmt.Errorf("instance %s is not running", name)

	}
	cmd := exec.Command("kill", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return err
	}
	inst.Instances[name].PID = 0
	return nil
}

func (inst *Controller) RestartInstance(name string) error {
	err := inst.StopInstance(name)
	if err != nil {
		return err
	}
	return inst.StartInstance(name)
}

func (inst *Controller) GetInstanceStatus(name string) string {
	port := inst.Instances[name].Port
	pid, err := getPID(port)
	if err != nil {
		return "PID not found"
	}
	return fmt.Sprintf("ruuning with PID: %d", pid)
}

func (inst *Controller) GetAllInstances() []*Instance {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()

	instances := make([]*Instance, 0, len(inst.Instances))
	for _, instance := range inst.Instances {
		pid, err := getPID(instance.Port)
		if err != nil {
			instance.Error = "instance isn't running"
		}
		instance.PID = pid
		instances = append(instances, instance)
	}
	return instances
}

func getPID(port int) (int, error) {
	cmd := exec.Command("fuser", "-n", "tcp", fmt.Sprint(port))
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	pidStr := strings.TrimSpace(string(out))
	return strconv.Atoi(pidStr)
}

func (inst *Controller) CreateInstance(c *gin.Context) {
	var instance Instance
	if err := c.BindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := inst.AddInstance(instance.Name, instance.Port)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = inst.SaveToFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance created successfully"})
}

func (inst *Controller) StartInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := inst.StartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance started successfully"})
}

func (inst *Controller) StopInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := inst.StopInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance stopped successfully"})
}

func (inst *Controller) RestartInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := inst.RestartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance restarted successfully"})
}

func (inst *Controller) DeleteInstanceHandler(c *gin.Context) {
	name := c.Query("name")
	err := inst.DeleteInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = inst.SaveToFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Instance deleted successfully"})
}

func (inst *Controller) GetInstanceStatusHandler(c *gin.Context) {
	name := c.Query("name")
	status := inst.GetInstanceStatus(name)
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (inst *Controller) GetAllInstancesHandler(c *gin.Context) {
	instances := inst.GetAllInstances()
	c.JSON(http.StatusOK, gin.H{"instances": instances})
}

func (inst *Controller) ReadYAMLFile(c *gin.Context) {
	yamlFile, err := os.ReadFile(DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/yaml", yamlFile)
}

func (inst *Controller) GetPIDByPortHandler(c *gin.Context) {
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

func (inst *Controller) ProxyHandler(c *gin.Context) {
	name := c.GetHeader("Instance-Name")
	instance, exists := inst.Instances[name]
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
