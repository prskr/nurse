package check

import (
	"golang.org/x/sync/errgroup"

	"code.1533b4dc0.de/prskr/nurse/config"
	"code.1533b4dc0.de/prskr/nurse/grammar"
)

var _ SystemChecker = Collection(nil)

type Collection []SystemChecker

func (Collection) UnmarshalCheck(grammar.Check, config.ServerLookup) error {
	panic("unmarshalling is not supported for a collection")
}

func (c Collection) Execute(ctx Context) error {
	grp, grpCtx := errgroup.WithContext(ctx)

	chkCtx := ctx.WithParent(grpCtx)

	for i := range c {
		chk := c[i]
		grp.Go(func() error {
			return chk.Execute(chkCtx)
		})
	}

	return grp.Wait()
}
