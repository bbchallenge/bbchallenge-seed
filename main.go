package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const DEBUG = 1

const R = 0
const L = 1

var nbStates byte

var mutexMetrics sync.Mutex
var nbMachineSeen int
var maxNbSteps int
var maxSpace int
var maxNbGoRoutines int

func printTM(tm TM) {
	for i := 0; i < int(nbStates); i += 1 {
		for j := 0; j <= 1; j += 1 {
			fmt.Printf("%d%d%d ", tm[6*i+3*j], tm[6*i+3*j+1], tm[6*i+3*j+2])
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

var start time.Time

// Invariant: tm's transition (state, read) is not defined
func search(tm TM, state byte, read byte,
	previous_steps_count int, previous_space_count int,
	slow_down_init int, slow_down int) {

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

				haltStatus, after_state, after_read, steps_count, space_count := simulate(newTm)

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
							search(newTm, after_state, after_read, steps_count, space_count, slow_down_init, slow_down_init)
							wg.Done()
						}()
					} else {
						search(newTm, after_state, after_read, steps_count, space_count,
							slow_down_init, slow_down-1)
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
	fmt.Println(nbMachineSeen, maxNbSteps, maxSpace, maxNbGoRoutines, time.Since(start), float64(nbMachineSeen)/time.Since(start).Seconds())
	mutexMetrics.Unlock()

}

func main() {
	start = time.Now()

	// +---+-----+-----+
	// | - |  0  |  1  |
	// +---+-----+-----+
	// | A | 1RB | 1LC |
	// | B | 1RC | 1RB |
	// | C | 1RD | 0LE |
	// | D | 1LA | 1LD |
	// | E | 1RH | 0LA |
	// +---+-----+-----+
	// bb5_winner := TM{
	// 	1, R, 2, 1, L, 3,
	// 	1, R, 3, 1, R, 2,
	// 	1, R, 4, 0, L, 5,
	// 	1, L, 1, 1, L, 4,
	// 	1, R, 6, 0, L, 1}

	// halt_status, curr_state, read, steps_count, space_count := simulate(bb5_winner)

	// fmt.Println(halt_status, curr_state, read,
	// 	steps_count, space_count)

	// Initial transition is 1RB (w.l.o.g)
	kick_start := TM{
		1, R, 2, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}

	nbStates = 5

	search(kick_start, 2, 0, 1, 1, 2, 2)

	fmt.Println(nbMachineSeen, maxNbSteps, maxSpace)
	fmt.Println("Max Go Routines:", maxNbGoRoutines)
}
