package model

type DeviceInfo struct {
	GlobalUUID  string `json:"globalUUID"`
	ClientId    string `json:"clientId"`
	ClientName  string `json:"clientName"`
	SiteId      string `json:"siteId"`
	SiteName    string `json:"siteName"`
	DeviceId    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	SiteAddress string `json:"siteAddress"`
	SiteCity    string `json:"siteCity"`
	SiteState   string `json:"siteState"`
	SiteZip     string `json:"siteZip"`
	SiteCountry string `json:"siteCountry"`
	SiteLat     string `json:"siteLat"`
	SiteLon     string `json:"siteLon"`
	TimeZone    string `json:"timeZone"`
	CreatedOn   string `json:"createdOn"`
	UpdatedOn   string `json:"updatedOn"`
}

type DeviceInfoFirstRecord struct {
	DeviceInfo DeviceInfo `json:"1"`
}

type DeviceInfoDefault struct {
	DeviceInfoFirstRecord DeviceInfoFirstRecord `json:"_default"`
}
