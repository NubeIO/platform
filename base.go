package platform

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const DB = "db.yaml"
const BinaryName = "main"

type Instance struct {
	Name string `yaml:"name"`
	Port int    `yaml:"port"`
	PID  int    `yaml:"pid"`
}

type InstanceManager struct {
	Instances map[string]*Instance
	Lock      sync.Mutex
}

func (im *InstanceManager) LoadFromFile(filePath string) error {
	im.Lock.Lock()
	defer im.Lock.Unlock()

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	im.Instances = make(map[string]*Instance)
	err = yaml.Unmarshal(yamlFile, &im.Instances)
	if err != nil {
		return err
	}

	return nil
}

func (im *InstanceManager) SaveToFile() error {
	im.Lock.Lock()
	defer im.Lock.Unlock()

	yamlData, err := yaml.Marshal(&im.Instances)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(DB, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (im *InstanceManager) AddInstance(name string, port int) error {
	im.Lock.Lock()
	defer im.Lock.Unlock()

	_, exists := im.Instances[name]
	if exists {
		return fmt.Errorf("instance with name %s already exists", name)
	}

	for _, instance := range im.Instances {
		if instance.Port == port {
			return fmt.Errorf("instance with port %d already exists", port)
		}
	}

	im.Instances[name] = &Instance{Name: name, Port: port}
	err := im.StartInstance(name)
	if err != nil {
		return err
	}
	return nil
}

func (im *InstanceManager) DeleteInstance(name string) error {
	im.Lock.Lock()
	defer im.Lock.Unlock()

	_, exists := im.Instances[name]
	if !exists {
		return fmt.Errorf("instance with name %s not found", name)
	}

	delete(im.Instances, name)
	err := im.StopInstance(name)
	if err != nil {
		return err
	}
	return nil
}

func (im *InstanceManager) StartInstance(name string) error {
	cmd := exec.Command(fmt.Sprintf("./%s", BinaryName), fmt.Sprintf("-p=%d", im.Instances[name].Port))
	err := cmd.Start()
	if err != nil {
		return err
	}
	im.Instances[name].PID = cmd.Process.Pid
	return nil
}

func (im *InstanceManager) StopInstance(name string) error {
	port := im.Instances[name].Port
	pid, err := getPID(port)
	if err != nil {
		return fmt.Errorf("instance %s is not running", name)

	}
	cmd := exec.Command("kill", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return err
	}
	im.Instances[name].PID = 0
	return nil
}

func (im *InstanceManager) RestartInstance(name string) error {
	err := im.StopInstance(name)
	if err != nil {
		return err
	}
	return im.StartInstance(name)
}

func (im *InstanceManager) GetInstanceStatus(name string) string {
	port := im.Instances[name].Port
	pid, err := getPID(port)
	if err != nil {
		return "PID not found"
	}
	cmd := exec.Command("kill", "-0", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		return "stopped"
	}
	return "running"
}

func (im *InstanceManager) GetAllInstances() []*Instance {
	im.Lock.Lock()
	defer im.Lock.Unlock()

	instances := make([]*Instance, 0, len(im.Instances))
	for _, instance := range im.Instances {
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
