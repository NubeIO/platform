package dto

type DiskUsage struct {
	Size      string `json:"size"`
	Used      string `json:"used"`
	Available string `json:"available"`
	Usage     string `json:"usage"`
}

type Disk struct {
	FileSystem string    `json:"fileSystem"`
	Type       string    `json:"type"`
	MountedOn  string    `json:"mountedOn"`
	Usage      DiskUsage `json:"usage"`
}
