package cmd

import (
	"fmt"
	"log/slog"

	"code.icb4dc0.de/prskr/nurse/check"
	"code.icb4dc0.de/prskr/nurse/grammar"
	"github.com/urfave/cli/v2"
)

type executor struct {
	*app
}

func (a *executor) ExecChecks(ctx *cli.Context) error {
	parser, err := grammar.NewParser[grammar.Script]()
	if err != nil {
		return fmt.Errorf("failed to create parser: %w", err)
	}

	var checks []grammar.Check

	for i := 0; i < ctx.NArg(); i++ {
		var s *grammar.Script
		if s, err = parser.Parse(ctx.Args().Get(i)); err != nil {
			return fmt.Errorf("failed to parse checks: %w", err)
		} else {
			checks = append(checks, s.Checks...)
		}
	}

	checker, err := check.CheckForScript(checks, a.registry, a.lookup)
	if err != nil {
		return fmt.Errorf("failed to compile checks: %w", err)
	}

	checkCtx, cancel := check.AttemptsContext(ctx.Context, a.nurseInstance.CheckAttempts, a.nurseInstance.CheckTimeout)
	defer cancel()

	if err := checker.Execute(checkCtx); err != nil {
		return err
	}

	slog.Default().Info("Successfully executed checks")

	return nil
}
