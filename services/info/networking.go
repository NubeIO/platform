package info

import (
	"errors"
	"fmt"
	"github.com/NubeIO/platform/model"
	"os/exec"
)

type NetworkingBody struct {
	PortName string `json:"port_name"`
}

func (inst *System) RestartNetworking() (*model.Message, error) {
	cmd := exec.Command("systemctl", "restart", "info.service")
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, false)
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Message: "restarted ok",
	}, err
}

func (inst *System) InterfaceUpDown(port NetworkingBody) (*model.Message, error) {
	_, err := inst.interfaceUpDown(port.PortName, false) // set down
	if err != nil {
		return nil, err
	}
	_, err = inst.interfaceUpDown(port.PortName, true) // set up
	if err != nil {
		return nil, err
	}
	return &model.Message{fmt.Sprintf("restart network: %s", port.PortName)}, nil

}

func (inst *System) InterfaceUp(port NetworkingBody) (*model.Message, error) {
	return inst.interfaceUpDown(port.PortName, true)
}

func (inst *System) InterfaceDown(port NetworkingBody) (*model.Message, error) {
	return inst.interfaceUpDown(port.PortName, false)
}

// interfaceUpDown ifconfig enp4s0 up
func (inst *System) interfaceUpDown(port string, up bool) (*model.Message, error) {
	if !portExists(port) {
		return nil, errors.New(fmt.Sprintf("port %s was not found", port))
	}
	cmd := exec.Command("ifconfig", port, "down")
	msg := "disabled"
	if up {
		msg = "enabled"
		cmd = exec.Command("ifconfig", port, "up")
	}
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Message: fmt.Sprintf("port %s is now %s", port, msg),
	}, err
}

func portExists(port string) bool {
	names, err := nets.GetInterfacesNames()
	if err != nil {
		return false
	}
	for _, s := range names.Names {
		if port == s {
			return true
		}
	}
	return false
}
