package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const defaultDirPermissions = 0o755

var (
	GoSourceFiles      []string
	GeneratedMockFiles []string
	WorkingDir         string
	OutDir             string
	dirsToIgnore       = []string{
		".git",
		"magefiles",
		".concourse",
		".run",
		".task",
	}
)

func init() {
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		WorkingDir = wd
	}

	OutDir = filepath.Join(WorkingDir, "out")

	if err := os.MkdirAll(OutDir, defaultDirPermissions); err != nil {
		panic(err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, nil)))

	if err := initSourceFiles(); err != nil {
		panic(err)
	}

	slog.Info("Completed initialization")
}

func initSourceFiles() error {
	return filepath.WalkDir(WorkingDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if slices.Contains(dirsToIgnore, filepath.Base(path)) {
				return fs.SkipDir
			}
			return nil
		}

		_, ext, found := strings.Cut(filepath.Base(path), ".")
		if !found {
			return nil
		}

		switch ext {
		case "mock.go":
			GeneratedMockFiles = append(GeneratedMockFiles, path)
		case "go":
			GoSourceFiles = append(GoSourceFiles, path)
		}

		return nil
	})
}
