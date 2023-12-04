package main

import (
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/magefile/mage/sh"
)

var (
	GoReleaser = sh.RunCmd("goreleaser")
	GoInstall  = sh.RunCmd("go", "install")
	GoBuild    = sh.RunCmd("go", "build")
)

func ensureGoTool(toolName, importPath, version string) error {
	return checkForTool(toolName, func() error {
		toolToInstall := fmt.Sprintf("%s@%s", importPath, version)
		slog.Info("Installing Go tool", slog.String("toolToInstall", toolToInstall))
		return GoInstall(toolToInstall)
	})
}

func checkForTool(toolName string, fallbackAction func() error) error {
	if _, err := exec.LookPath(toolName); err != nil {
		slog.Warn("tool is missing", slog.String("toolName", toolName))
		return fallbackAction()
	}

	return nil
}
