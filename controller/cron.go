package controller

import (
	"errors"
	"fmt"
	"github.com/NubeIO/platform/dto"
	"github.com/NubeIO/platform/services/crontab"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"strings"
)

func (inst *Controller) GetRestartJob(c *gin.Context) {
	restartJobs := crontab.List()
	if restartJobs == nil {
		responseHandler([]dto.RestartJob{}, nil, c)
		return
	}
	responseHandler(restartJobs, nil, c)
}

func (inst *Controller) UpdateRestartJob(c *gin.Context) {
	body, err := getBodyRestartJob(c)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	err = validateCornExpression(body.Expression)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	err = crontab.Put(body)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(body, nil, c)
}

func (inst *Controller) DeleteRestartJob(c *gin.Context) {
	unit := c.Param("unit")
	err := crontab.Delete(unit)
	if err != nil {
		responseHandler(nil, err, c)
		return
	}
	responseHandler(dto.Message{Message: fmt.Sprintf("deleted %s restart job successfully", unit)}, nil, c)
}

func getBodyRestartJob(ctx *gin.Context) (dto *dto.RestartJob, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func validateCornExpression(exp string) error {
	expFields := strings.Fields(exp)
	if len(expFields) != 5 {
		return errors.New("invalid expression")
	}
	_, err := cron.ParseStandard(exp)
	if err != nil {
		return errors.New("invalid expression")
	}
	return nil
}
