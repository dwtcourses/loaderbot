/*
 * // Copyright 2020 Insolar Network Ltd.
 * // All rights reserved.
 * // This material is licensed under the Insolar License version 1.0,
 * // available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.
 */

package loaderbot

import (
	"context"
	"sync"
	"testing"
	"time"
)

func DefaultRunnerCfg() *RunnerConfig {
	return &RunnerConfig{
		Name:            "test_runner",
		Attackers:       1,
		AttackerTimeout: 1,
		StartRPS:        20,
		StepDurationSec: 5,
		StepRPS:         5,
		TestTimeSec:     60,
		ReportOptions: &ReportOptions{
			CSV: false,
			PNG: false,
		},
	}
}

func TestAttackSuccess(t *testing.T) {
	r := NewRunner(DefaultRunnerCfg(), &ControlAttackerMock{}, nil)
	r.controlled.Sleep = 10
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Cfg.TestTimeSec)*time.Second)
	r.TimeoutCtx = ctx
	r.CancelFunc = cancel

	// sync
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go attack(r.attackers[0], r, wg)
	wg.Wait()

	r.next <- attackToken{
		Step: 1,
		Tick: 1,
	}
	res := <-r.results
	if got, want := res.DoResult.Error, ""; got != want {
		t.Fatalf("got %v want %v", got, want)
	}
	if got, want := int(res.Elapsed), int(r.controlled.Sleep); got < want {
		t.Fatalf("got %v want >= %v", got, want)
	}
	// async
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go asyncAttack(r.attackers[0], r, wg2)
	wg2.Wait()

	r.next <- attackToken{
		Step: 1,
		Tick: 1,
	}
	res2 := <-r.results
	if got, want := res2.DoResult.Error, ""; got != want {
		t.Fatalf("got %v want %v", got, want)
	}
	if got, want := int(res.Elapsed), int(r.controlled.Sleep); got < want {
		t.Fatalf("got %v want >= %v", got, want)
	}
}

func TestAttackTimeout(t *testing.T) {
	r := NewRunner(DefaultRunnerCfg(), &ControlAttackerMock{}, nil)
	r.controlled.Sleep = 2000
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Cfg.TestTimeSec)*time.Second)
	r.TimeoutCtx = ctx
	r.CancelFunc = cancel

	// sync
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go attack(r.attackers[0], r, wg)
	wg.Wait()

	r.next <- attackToken{
		Step: 1,
		Tick: 1,
	}
	res := <-r.results
	if got, want := res.DoResult.Error, errAttackDoTimedOut; got != want {
		t.Fatalf("got %v want %v", got, want)
	}

	// async
	wg2 := &sync.WaitGroup{}
	wg2.Add(1)
	go attack(r.attackers[0], r, wg2)
	wg2.Wait()

	r.next <- attackToken{
		Step: 1,
		Tick: 1,
	}
	res2 := <-r.results
	if got, want := res2.DoResult.Error, errAttackDoTimedOut; got != want {
		t.Fatalf("got %v want %v", got, want)
	}
}
