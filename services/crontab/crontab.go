package crontab

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/NubeIO/lib-utils-go/nstring"
	"github.com/NubeIO/platform/dto"
	"os/exec"
	"strings"
	"time"
)

func List() []*dto.RestartJob {
	restartJobs := make([]*dto.RestartJob, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	output, _, _, err := execute(ctx, "crontab", "-l")
	if err != nil {
		return restartJobs
	}
	outputLines := strings.Split(output, "\n")
	for _, l := range outputLines {
		l = strings.Trim(l, " ")
		if !strings.HasPrefix(l, "#") && strings.Contains(l, "systemctl restart rubix") {
			parts := strings.SplitN(l, " systemctl restart ", 2)
			if len(parts) == 2 {
				restartJob := &dto.RestartJob{}
				if !strings.HasSuffix(parts[1], ".service") {
					restartJob.Unit = strings.Trim(fmt.Sprintf("%s.service", parts[1]), " ")
				} else {
					restartJob.Unit = strings.Trim(parts[1], " ")
				}
				restartJob.Expression = strings.Trim(parts[0], " ")
				restartJobs = append(restartJobs, restartJob)
			}
		}
	}
	return restartJobs
}

func Get(unit string) *string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	unitWithoutExtension := strings.Replace(unit, ".service", "", 1)
	output, _, _, err := execute(
		ctx,
		"bash",
		"-c",
		fmt.Sprintf("crontab -l | grep -E '%s$|%s'", unitWithoutExtension, unit),
	)
	if err != nil {
		return nil
	}
	output = strings.Trim(output, " ")
	if !strings.HasPrefix(output, "#") && strings.Contains(output, "systemctl restart rubix") {
		parts := strings.SplitN(output, " systemctl restart ", 2)
		if len(parts) == 2 {
			return nstring.New(strings.Trim(parts[0], " "))
		}
	}
	return nil
}

func Put(restartJob *dto.RestartJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	unit := restartJob.Unit
	if !strings.HasPrefix(unit, "rubix-") {
		return errors.New(fmt.Sprintf("correct the unit %s", unit))
	}
	if !strings.HasSuffix(unit, ".service") {
		unit = fmt.Sprintf("%s.service", unit)
	}

	unitWithoutExtension := strings.Replace(unit, ".service", "", 1)
	_, _, _, err := execute(
		ctx,
		"bash",
		"-c",
		fmt.Sprintf("(crontab -l | grep -v '%s$' | grep -v '%s') | crontab -", unitWithoutExtension, unit),
	)
	if err != nil {
		return err
	}

	_, _, _, err = execute(
		ctx,
		"bash",
		"-c",
		fmt.Sprintf("(crontab -l ; echo '%s systemctl restart %s') | crontab -", restartJob.Expression, unit),
	)

	if err != nil {
		return err
	}

	return nil
}

func Delete(unit string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !strings.HasPrefix(unit, "rubix-") {
		return errors.New(fmt.Sprintf("correct the unit %s", unit))
	}
	if !strings.HasSuffix(unit, ".service") {
		unit = fmt.Sprintf("%s.service", unit)
	}

	unitWithoutExtension := strings.Replace(unit, ".service", "", 1)
	_, _, _, err := execute(
		ctx,
		"bash",
		"-c",
		fmt.Sprintf("(crontab -l | grep -v '%s$' | grep -v '%s') | crontab -", unitWithoutExtension, unit),
	)
	if err != nil {
		return err
	}
	return nil
}

func execute(ctx context.Context, command string, args ...string) (string, string, int, error) {
	var (
		err      error
		stderr   bytes.Buffer
		stdout   bytes.Buffer
		code     int
		output   string
		warnings string
	)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	output = stdout.String()
	warnings = stderr.String()
	code = cmd.ProcessState.ExitCode()
	return output, warnings, code, err
}
