package info

import (
	"net"
	"strings"
)

const (
	active                    = "Active: active (running)"
	dead                      = "Active: inactive (dead)"
	statusOk                  = "Status: Initialization Sequence Completed"
	processRestarting         = "received, process restarting"
	couldNotDetermineProtocol = "Could not determine IPv4/IPv6 protocol"
	cannotResolveHostAddress  = "Cannot resolve host address"
)

const openVPNServiceName = "openvpn@client"

type VPNStatus struct {
	Active                    bool     `json:"active"`
	Dead                      bool     `json:"dead"`
	StatusOk                  bool     `json:"status_ok"`
	ProcessRestarting         bool     `json:"processRestarting"`
	CouldNotDetermineProtocol bool     `json:"couldNotDetermineProtocol"`
	CannotResolveHostAddress  bool     `json:"cannotResolveHostAddress"`
	Messages                  []string `json:"messages"`
	Ip                        string   `json:"ip"`
}

func (inst *System) OpenVPNStatus() (*VPNStatus, error) {
	data, err := inst.systemctl.Status(openVPNServiceName)
	if err != nil {
		return nil, err
	}
	var message []string
	var isActive bool
	var isDead bool
	var ok bool
	var isProcessRestarting bool
	var isCouldNotDetermineProtocol bool
	var isCannotResolveHostAddress bool

	if strings.Contains(data, active) {
		isActive = true
		message = append(message, active)
	}
	if strings.Contains(data, dead) {
		isDead = true
		message = append(message, dead)
	}
	if strings.Contains(data, statusOk) {
		ok = true
		message = append(message, statusOk)
	}
	if strings.Contains(data, processRestarting) {
		isProcessRestarting = true
		message = append(message, processRestarting)
	}
	if strings.Contains(data, couldNotDetermineProtocol) {
		isCouldNotDetermineProtocol = true
		message = append(message, couldNotDetermineProtocol)
	}
	if strings.Contains(data, cannotResolveHostAddress) {
		isCannotResolveHostAddress = true
		message = append(message, cannotResolveHostAddress)
	}

	ip := "ip not found on interface tun0"
	getIp := getInternalIP()
	if getIp != "" {
		ip = getIp
	}

	out := &VPNStatus{
		Active:                    isActive,
		Dead:                      isDead,
		StatusOk:                  ok,
		ProcessRestarting:         isProcessRestarting,
		CouldNotDetermineProtocol: isCouldNotDetermineProtocol,
		CannotResolveHostAddress:  isCannotResolveHostAddress,
		Messages:                  message,
		Ip:                        ip,
	}
	return out, nil

}

func getInternalIP() string {
	itf, _ := net.InterfaceByName("tun0")
	item, _ := itf.Addrs()
	var ip net.IP
	for _, addr := range item {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil {
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String()
	} else {
		return ""
	}
}
