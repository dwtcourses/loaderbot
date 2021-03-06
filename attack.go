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
	"time"
)

// Attack must be implemented by a service client.
type Attack interface {
	// Setup should establish the connection to the service
	// It may want to access the Config of the Runner.
	Setup(c RunnerConfig) error
	// Do performs one request and is executed in a separate goroutine.
	// The context is used to CancelFunc the request on timeout.
	Do(ctx context.Context) DoResult
	// Teardown can be used to close the connection to the service
	Teardown() error
	// Clone should return a fresh new Attack
	// Make sure the new Attack has values for shared struct fields initialized at Setup.
	Clone(r *Runner) Attack
}

// attack receives schedule signal and attacks target calling Do() method, returning AttackResult with timings
func attack(a Attack, r *Runner, wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case <-r.TimeoutCtx.Done():
			return
		case nextMsg := <-r.next:
			requestCtx, requestCtxCancel := context.WithTimeout(context.Background(), time.Duration(r.Cfg.AttackerTimeout)*time.Second)

			tStart := time.Now()

			done := make(chan DoResult)

			var doResult DoResult

			go func() {
				done <- a.Do(requestCtx)
			}()
			// either get the result from the attacker or from the timeout
			select {
			case <-r.TimeoutCtx.Done():
				requestCtxCancel()
				return
			case <-requestCtx.Done():
				doResult = DoResult{
					RequestLabel: r.Name,
					Error:        errAttackDoTimedOut,
				}
			case doResult = <-done:
			}

			tEnd := time.Now()

			atkResult := AttackResult{
				AttackToken: nextMsg,
				Begin:       tStart,
				End:         tEnd,
				Elapsed:     tEnd.Sub(tStart),
				DoResult:    doResult,
			}
			requestCtxCancel()
			r.results <- atkResult
		}
	}
}

// asyncAttack receives schedule signal and attacks target calling Do() method asynchronously, returning AttackResult with timings
func asyncAttack(a Attack, r *Runner, wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case <-r.TimeoutCtx.Done():
			return
		case nextMsg := <-r.next:
			requestCtx, requestCtxCancel := context.WithTimeout(context.Background(), time.Duration(r.Cfg.AttackerTimeout)*time.Second)

			tStart := time.Now()

			done := make(chan DoResult)

			var doResult DoResult

			go func() {
				done <- a.Do(requestCtx)
			}()

			go func() {
				// either get the result from the attacker or from the timeout
				select {
				case <-r.TimeoutCtx.Done():
					requestCtxCancel()
					return
				case <-requestCtx.Done():
					doResult = DoResult{
						RequestLabel: r.Name,
						Error:        errAttackDoTimedOut,
					}
				case doResult = <-done:
				}

				tEnd := time.Now()

				atkResult := AttackResult{
					AttackToken: nextMsg,
					Begin:       tStart,
					End:         tEnd,
					Elapsed:     tEnd.Sub(tStart),
					DoResult:    doResult,
				}
				requestCtxCancel()
				r.results <- atkResult
			}()
		}
	}
}
