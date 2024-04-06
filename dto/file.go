package dto

type FileExistence struct {
	File   string `json:"file"`
	Exists bool   `json:"exists"`
}

type FileUploadResponse struct {
	Destination string `json:"destination"`
	File        string `json:"file"`
	Size        string `json:"size"`
	UploadTime  string `json:"uploadTime"`
}

type WriteFile struct {
	FilePath     string     `json:"path"`
	Body         DynamicMap `json:"body"`
	BodyAsString string     `json:"bodyAsString"`
}

type WriteFileData struct {
	Data string `json:"data"`
}
