package loaderbot

import (
	"sync/atomic"
	"time"
)

func serviceErrorAfter(se chan bool, t time.Duration) {
	go func() {
		time.Sleep(t)
		se <- true
	}()
}

type ControllableConfig struct {
	R               *Runner
	ControlChan     chan bool
	AttackerLatency time.Duration
	AttackersAmount int
}

func withControllableAttackers(cfg ControllableConfig) {
	attackers := make([]Attack, 0)
	for i := 0; i < cfg.AttackersAmount; i++ {
		attackers = append(attackers, NewControlMockAttacker(i, cfg.ControlChan, cfg.R))
	}
	cfg.R.attackers = attackers
}

type ServiceLatencyChangeConfig struct {
	R             *Runner
	Interval      time.Duration
	LatencyStepMs uint64
	Times         int
	LatencyFlag   int
}

const (
	increaseLatency = iota
	decreaseLatency
)

func changeAttackersLatency(cfg ServiceLatencyChangeConfig) {
	go func() {
		for i := 0; i < cfg.Times; i++ {
			if cfg.LatencyFlag == increaseLatency {
				atomic.AddUint64(&cfg.R.controlled.Sleep, cfg.LatencyStepMs)
			}
			if cfg.LatencyFlag == decreaseLatency {
				atomic.AddUint64(&cfg.R.controlled.Sleep, -cfg.LatencyStepMs)
			}
			time.Sleep(cfg.Interval)
		}
		cfg.R.L.Infof("=== done changing latency ===")
		cfg.R.L.Infof("=== keeping latency constant for new attackers ===")
	}()
}