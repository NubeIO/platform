package dto

type Plugin struct {
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	Version   string `json:"version,omitempty"`
	Extension string `json:"extension"`
}

// PluginExternal holds information about a plugin instance for one user.
type PluginExternal struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	Author       string   `json:"author,omitempty" form:"author" query:"author"`
	Website      string   `json:"website,omitempty" form:"website" query:"website"`
	License      string   `json:"license,omitempty" form:"license" query:"license"`
	HasNetwork   bool     `json:"has_network"`
	Capabilities []string `json:"capabilities,omitempty"`
}
