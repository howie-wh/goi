package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

// ExecuteCmd ...
func ExecuteCmd(cmdStr string, path string) error {
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Infof("path :%v", path)
	log.Infof("cmd is: %v", cmd.Args)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("fail to execute: %v, err: %w", cmd.Args, err)
	}
	return nil
}
