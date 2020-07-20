package loaderbot

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	sigs = make(chan os.Signal, 1)
)

func (r *Runner) handleShutdownSignal() {
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		r.L.Infof("exit signal received, exiting")
		if r.Cfg.GoroutinesDump {
			buf := make([]byte, 1<<20)
			stacklen := runtime.Stack(buf, true)
			r.L.Infof("=== received SIGTERM ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stacklen])
		}
		os.Exit(1)
	}()
}

func NewImmediateTicker(repeat time.Duration) *time.Ticker {
	ticker := time.NewTicker(repeat)
	oc := ticker.C
	nc := make(chan time.Time, 1)
	go func() {
		nc <- time.Now()
		for tm := range oc {
			nc <- tm
		}
	}()
	ticker.C = nc
	return ticker
}

func MaxRPS(array []float64) float64 {
	if len(array) == 0 {
		return 1
	}
	var max = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
	}
	return max
}