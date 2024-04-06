package info

import (
	"fmt"
	"github.com/NubeIO/lib-date/datelib"
	"github.com/NubeIO/platform/model"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

type DateBody struct {
	DateTime string `json:"dateTime"`
	TimeZone string `json:"timezone"`
}

func (inst *System) GetTimezoneList() ([]string, error) {
	return inst.datectl.GetTimezoneList()
}

func (inst *System) UpdateTimezone(body DateBody) (*model.Message, error) {
	err := inst.datectl.UpdateTimezone(body.TimeZone)
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Message: fmt.Sprintf("updated to %s", body.TimeZone),
	}, nil
}

func (inst *System) SetSystemTime(dateTime DateBody) (*datelib.Time, error) {
	layout := "2006-01-02 15:04:05"
	// parse time
	t, err := time.Parse(layout, dateTime.DateTime)
	if err != nil {
		return nil, fmt.Errorf("could not parse date try 2006-01-02 15:04:05 %s", err)
	}
	log.Infof("set time to %s", t.String())
	timeString := fmt.Sprintf("%s", dateTime.DateTime)
	cmd := exec.Command("date", "-s", timeString)
	output, err := cmd.Output()
	cleanCommand(string(output), cmd, err, debug)
	if err != nil {
		return nil, err
	}
	dateLib := datelib.Date{}
	return dateLib.SystemTime(), nil
}

func (inst *System) NTPEnable() (*model.Message, error) {
	msg, err := inst.datectl.NTPEnable()
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Message: msg.Message,
	}, nil
}

func (inst *System) NTPDisable() (*model.Message, error) {
	msg, err := inst.datectl.NTPDisable()
	if err != nil {
		return nil, err
	}
	return &model.Message{
		Message: msg.Message,
	}, nil
}
