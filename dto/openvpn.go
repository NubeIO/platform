package dto

type OpenVPNClient struct {
	VirtualIP      string `json:"virtualIP"`
	ReceivedBytes  int    `json:"receivedBytes"`
	SentBytes      int    `json:"sentBytes"`
	ConnectedSince string `json:"connectedSince"`
}

type OpenVPNBody struct {
	Name string `json:"name"`
}

type OpenVPNConfig struct {
	Data string `json:"data"`
}
