package rubixregistry

import (
	"encoding/json"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/times/utilstime"
	"github.com/NubeIO/platform/constants"
	"github.com/NubeIO/platform/dto"
	"github.com/NubeIO/platform/utils/crontab"
	"os"
)

func (inst *RubixRegistry) GetDeviceInfo() (*dto.DeviceInfo, error) {
	data, err := os.ReadFile(inst.GlobalUUIDFile)
	if err != nil {
		return nil, err
	}
	info, err := inst.GetProductInfo()
	if err != nil {
		return nil, err
	}
	var model string
	var version string
	if info != nil {
		model = info.Type
		version = info.Version
	}
	timezone, _ := utilstime.GetHardwareTZ()
	out := &dto.DeviceInfo{
		GlobalUUID: string(data),
		Version:    version,
		Type:       model,
		Timezone:   timezone,
		ROS: dto.ROSInfo{
			Version:           inst.GetInstalledAppVersion(constants.RubixOs),
			RestartExpression: crontab.Get("nubeio-rubix-os.service"),
		},
	}
	return out, nil
}

func (inst *RubixRegistry) GetLegacyDeviceInfo() (*dto.DeviceInfo, error) { // TODO: remove after migration done
	data, err := os.ReadFile(inst.LegacyDeviceInfoFile)
	if err != nil {
		return nil, err
	}
	deviceInfoDefault := dto.DeviceInfoDefault{}
	err = json.Unmarshal(data, &deviceInfoDefault)
	if err != nil {
		return nil, err
	}
	return &deviceInfoDefault.DeviceInfoFirstRecord.DeviceInfo, nil
}

func (inst *RubixRegistry) GetProductInfo() (*dto.ProductInfo, error) {
	data, err := os.ReadFile(inst.ProductInfoPath)
	if err != nil {
		return &dto.ProductInfo{}, nil
	}
	productInfoDefault := dto.ProductInfo{}
	err = json.Unmarshal(data, &productInfoDefault)
	if err != nil {
		return &dto.ProductInfo{}, nil
	}
	return &productInfoDefault, nil
}
