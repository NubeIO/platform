package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
)

const DB = "db.yaml"

type Instance struct {
	Name                        string   `json:"name" yaml:"name"`
	Repo                        string   `json:"repo" yaml:"repo"`
	Description                 string   `json:"description" yaml:"description"`
	Port                        int      `json:"port,omitempty" yaml:"port,omitempty"`
	Transport                   string   `json:"transport" yaml:"transport"`
	ExecStart                   string   `json:"execStart" yaml:"exec_start"`
	AttachWorkingDirOnExecStart bool     `json:"attachWorkingDirOnExecStart" yaml:"attach_working_dir_on_exec_start"`
	EnvironmentVars             []string `json:"environmentVars" yaml:"environment_vars"`
	Products                    []string `json:"products" yaml:"products"`
	Arch                        []string `json:"arch" yaml:"arch"`
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

func (inst *Controller) AddInstance(instance *Instance) error {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()
	name := instance.Name
	_, exists := inst.Instances[instance.Name]
	if exists {
		return fmt.Errorf("instance with name %s already exists", name)
	}

	inst.Instances[name] = instance
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

	return nil
}

func (inst *Controller) StopInstance(name string) error {

	return nil
}

func (inst *Controller) RestartInstance(name string) error {
	err := inst.StopInstance(name)
	if err != nil {
		return err
	}
	return inst.StartInstance(name)
}

func (inst *Controller) GetInstance(name string) (*Instance, error) {
	instance := inst.Instances[name]
	if instance == nil {
		return nil, fmt.Errorf("not found: %s", name)
	}
	return instance, nil
}

func (inst *Controller) GetAllInstances() []*Instance {
	inst.Lock.Lock()
	defer inst.Lock.Unlock()
	instances := make([]*Instance, 0, len(inst.Instances))
	for _, instance := range inst.Instances {
		instances = append(instances, instance)
	}
	return instances
}

func (inst *Controller) CreateInstance(c *gin.Context) {
	var instance *Instance
	if err := c.BindJSON(&instance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if instance == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse json"})
		return
	}

	err := inst.AddInstance(instance)
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
	name := c.Param("name")
	err := inst.StartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance started successfully"})
}

func (inst *Controller) StopInstanceHandler(c *gin.Context) {
	name := c.Param("name")
	err := inst.StopInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance stopped successfully"})
}

func (inst *Controller) RestartInstanceHandler(c *gin.Context) {
	name := c.Param("name")
	err := inst.RestartInstance(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Instance restarted successfully"})
}

func (inst *Controller) DeleteInstanceHandler(c *gin.Context) {
	name := c.Param("name")
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

func (inst *Controller) GetInstancesHandler(c *gin.Context) {
	name := c.Param("name")
	instance, err := inst.GetInstance(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"instances": instance})
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
