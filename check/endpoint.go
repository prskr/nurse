package check

import (
	"code.icb4dc0.de/prskr/nurse/config"
	"code.icb4dc0.de/prskr/nurse/grammar"
)

func CheckForScript(script []grammar.Check, lkp ModuleLookup, srvLookup config.ServerLookup) (Collection, error) {
	compiledChecks := make([]SystemChecker, 0, len(script))

	for i := range script {
		rawChk := script[i]
		mod, err := lkp.Lookup(rawChk.Initiator.Module)
		if err != nil {
			return nil, err
		}

		compiledCheck, err := mod.Lookup(rawChk, srvLookup)
		if err != nil {
			return nil, err
		}

		compiledChecks = append(compiledChecks, compiledCheck)
	}

	return compiledChecks, nil
}
