package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type SimulationBackend byte

const (
	SIMULATION_GO SimulationBackend = iota
	SIMULATION_C
)

var mutexMetrics sync.Mutex
var nbMachineSeen int
var maxNbSteps int
var maxSpace int
var maxNbGoRoutines int

// Invariant: tm's transition (state, read) is not defined
func search(tm TM, state byte, read byte,
	previous_steps_count int, previous_space_count int,
	slow_down_init int, slow_down int, simulation_backend SimulationBackend) {

	// mutexMetrics.Lock()
	// printTM(tm)
	// mutexMetrics.Unlock()

	// Get the list of candidate target states
	// taking all states up to the first completely undefined
	// As in http://turbotm.de/~heiner/BB/mabu90.html#Enumeration
	var target_states [MAX_STATES]byte
	var undefinedTransitionCount byte

	for iState := 0; iState < int(nbStates); iState += 1 {
		target_states[iState] = byte(iState + 1)

		if tm[6*(iState)+3*0+2] == 0 {
			undefinedTransitionCount += 1
		}

		if tm[6*(iState)+3*1+2] == 0 {
			undefinedTransitionCount += 1
		}

		// The following allows to take the min of the completely undefined states
		// last condition very important to discard current state (which is about to not be completely undefined)
		if tm[6*(iState)+3*0+2] == 0 && tm[6*(iState)+3*1+2] == 0 && byte(iState+1) != state {
			break
		}
	}

	// Last transition
	if target_states[nbStates-1] != 0 && undefinedTransitionCount == 1 {
		return
	}

	var target_state byte

	var localNbMachineSeen int
	var localMaxNbSteps int
	var localMaxSpace int

	var wg sync.WaitGroup
	for _, target_state = range target_states {
		if target_state == 0 {
			break
		}

		var move byte
		for move = 0; move <= 1; move += 1 {

			var write byte
			for write = 0; write <= 1; write += 1 {

				var newTm TM = tm
				newTm[(state-1)*6+read*3] = write
				newTm[(state-1)*6+read*3+1] = move
				newTm[(state-1)*6+read*3+2] = target_state
				localNbMachineSeen += 1

				var haltStatus HaltStatus
				var after_state byte
				var after_read byte
				var steps_count int
				var space_count int

				switch simulation_backend {
				case SIMULATION_GO:
					haltStatus, after_state, after_read, steps_count, space_count = simulate(newTm)
					break
				case SIMULATION_C:
					haltStatus, after_state, after_read, steps_count, space_count = simulate_C_wrapper(newTm)
					break
				}

				if haltStatus == HALT {

					// mutexMetrics.Lock()
					// if steps_count >= maxNbSteps {
					// 	fmt.Println(steps_count)
					// 	printTM(newTm)
					// }
					// mutexMetrics.Unlock()

					localMaxNbSteps = MaxI(localMaxNbSteps, steps_count)
					localMaxSpace = MaxI(localMaxSpace, space_count)

					if slow_down == 0 {
						wg.Add(1)

						go func() {
							search(newTm, after_state, after_read, steps_count, space_count,
								slow_down_init, slow_down_init, simulation_backend)
							wg.Done()
						}()
					} else {
						search(newTm, after_state, after_read, steps_count, space_count,
							slow_down_init, slow_down-1, simulation_backend)
					}

				}

			}

		}

	}
	wg.Wait()

	mutexMetrics.Lock()
	nbMachineSeen += localNbMachineSeen

	maxNbSteps = MaxI(localMaxNbSteps, maxNbSteps)
	maxSpace = MaxI(localMaxSpace, maxSpace)
	maxNbGoRoutines = MaxI(maxNbGoRoutines, runtime.NumGoroutine())

	if VERBOSITY >= 2 {
		fmt.Println(nbMachineSeen, maxNbSteps, maxSpace, maxNbGoRoutines, time.Since(start), float64(nbMachineSeen)/time.Since(start).Seconds())
	}
	mutexMetrics.Unlock()

}
