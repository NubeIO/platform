package dto

type IPType struct {
	REST string `json:"rest"`
	UDP  string `json:"udp"`
	MQTT string `json:"mqttClient"`
}

// IPNetwork type ip based network
type IPNetwork struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Token    string `json:"token"`
	IPType
}

type NetworkTagForPostgresSync struct {
	HostUUID    string `json:"hostUUID"`
	NetworkUUID string `json:"networkUUID"`
	Tag         string `json:"tag"`
}

type NetworkMetaTagForPostgresSync struct {
	HostUUID    string `json:"hostUUID,omitempty"`
	NetworkUUID string `json:"networkUUID,omitempty"`
	Key         string `json:"key,omitempty"`
	Value       string `json:"value,omitempty"`
}
