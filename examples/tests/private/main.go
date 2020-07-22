/*
 * // Copyright 2020 Insolar Network Ltd.
 * // All rights reserved.
 * // This material is licensed under the Insolar License version 1.0,
 * // available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.
 */

package main

import (
	"fmt"

	"github.com/insolar/loaderbot"
	"github.com/insolar/loaderbot/examples/attackers"
)

func main() {
	// Private system, constant amount of attackers
	cfg := &loaderbot.RunnerConfig{
		TargetUrl:       "https://clients5.google.com/pagead/drt/dn/",
		Name:            "runner_1",
		SystemMode:      loaderbot.PrivateSystem,
		Attackers:       100,
		AttackerTimeout: 5,
		StartRPS:        100,
		StepDurationSec: 5,
		StepRPS:         10,
		TestTimeSec:     200,
	}
	lt := loaderbot.NewRunner(cfg, &attackers.AttackerExample{}, nil)
	maxRPS, _ := lt.Run()
	fmt.Printf("max rps: %.2f", maxRPS)
}