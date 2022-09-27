package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/sh"
	"go.uber.org/zap"
)

var (
	GoReleaser = sh.RunCmd("goreleaser")
	GoInstall  = sh.RunCmd("go", "install")
	GoBuild    = sh.RunCmd("go", "build")
)

func ensureGoTool(toolName, importPath, version string) error {
	return checkForTool(toolName, func() error {
		logger := zap.L()
		toolToInstall := fmt.Sprintf("%s@%s", importPath, version)
		logger.Info("Installing Go tool", zap.String("toolToInstall", toolToInstall))
		return GoInstall(toolToInstall)
	})
}

func checkForTool(toolName string, fallbackAction func() error) error {
	logger := zap.L()
	if _, err := exec.LookPath(toolName); err != nil {
		logger.Warn("tool is missing", zap.String("toolName", toolName))
		return fallbackAction()
	}

	return nil
}
