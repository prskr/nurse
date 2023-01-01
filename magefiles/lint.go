package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Lint(ctx context.Context) {
	mg.CtxDeps(ctx, Format, LintGo)
}

func LintGo(ctx context.Context) (err error) {
	return sh.RunV(
		"golangci-lint",
		"run",
		"-v",
		"--issues-exit-code=1",
	)
}
