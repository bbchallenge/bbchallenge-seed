// Here we define the TM enumeration algorithm
package bbchallenge

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"time"
)

type SimulationBackend byte

const (
	SIMULATION_GO SimulationBackend = iota
	SIMULATION_C
)

// Package parameters

var TimeStart time.Time = time.Now()

var UndecidedTimeLog io.Writer  // Logging UNDECIDED_TIME machines
var UndecidedSpaceLog io.Writer // Logging UNDECIDED_SPACE machines
var BBRecordLog io.Writer       // Logging BB and BB_space record holders

var Verbose bool
var LogFreq int64 = 30000000000 // 30 sec in ns

var ActivateFiltering bool = true

var SimulationLimitTime int = BB5
var SimulationLimitSpace int = BB5_SPACE

var SlowDownInit int = 2 // How many recursion will be done on the stack before calling go routines

// At the root of the TM tree are always 12 machines (independently of nbStates)
// Only 8 of them are interesting. We authorize the user to cut the computation
// of the entire tree in 1, 2, 4, or 8 pieces. Useful when you have several computers!
var TaskDivisor int = 1   // Can be either 1, 2, 4 or 8
var TaskDivisorMe int = 0 // Wich task to I do

// Package outputs

var mutexMetrics sync.Mutex
var NbMachineSeen int
var NbMachinePruned int
var NbHaltingMachines int
var NbNonHaltingMachines int
var NbUndecidedTime int
var NbUndecidedSpace int
var MaxNbSteps int
var MaxSpace int
var MaxNbGoRoutines int

// Logging related internal variables
var lastLogTime time.Time
var notFirstLog bool

// Invariant: tm's transition (state, read) is not defined
func Enumerate(nbStates byte, tm TM, state byte, read byte,
	previous_steps_count int, previous_space_count int,
	slow_down int, simulation_backend SimulationBackend) {

	// mutexMetrics.Lock()
	// printTM(tm)
	// mutexMetrics.Unlock()

	// Get the list of candidate target states
	// taking all states up to the first completely undefined
	// As in http://turbotm.de/~heiner/BB/mabu90.html#Enumeration
	var target_states [MAX_STATES]byte
	var definedTransitionCount byte
	var undefinedTransitionCount byte

	for iState := 0; iState < int(nbStates); iState += 1 {
		target_states[iState] = byte(iState + 1)

		if tm[6*(iState)+3*0+2] == 0 {
			undefinedTransitionCount += 1
		} else {
			definedTransitionCount += 1
		}

		if tm[6*(iState)+3*1+2] == 0 {
			undefinedTransitionCount += 1
		} else {
			definedTransitionCount += 1
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

	isRoot := definedTransitionCount == 1

	var target_state byte

	var localNbMachineSeen int
	var localNbMachinePruned int
	var localNbHalt int
	var localNbNoHalt int
	var localNbUndecidedTime int
	var localNbUndecidedSpace int
	var localMaxNbSteps int
	var localBestTimeHaltingMachine TM
	var localBestSpaceHaltingMachine TM
	var localMaxSpace int

	var loopIndex int

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

				if !isRoot && ActivateFiltering && pruneTM(nbStates, newTm, state, read) {
					localNbMachinePruned += 1
					continue
				}

				localNbMachineSeen += 1

				var haltStatus HaltStatus
				var after_state byte
				var after_read byte
				var steps_count int
				var space_count int

				switch simulation_backend {
				case SIMULATION_GO:
					haltStatus, after_state, after_read, steps_count, space_count = simulate(newTm, SimulationLimitTime, SimulationLimitSpace)
					break
				case SIMULATION_C:
					haltStatus, after_state, after_read, steps_count, space_count = simulate_C_wrapper(newTm, SimulationLimitTime, SimulationLimitSpace)
					break
				}

				switch haltStatus {
				case HALT:

					// Task Divisor
					if isRoot {

						if loopIndex/(8/TaskDivisor) != TaskDivisorMe {
							loopIndex += 1
							continue
						} else {
							loopIndex += 1
						}
					}

					if steps_count > localMaxNbSteps {
						localBestTimeHaltingMachine = newTm
					}

					if space_count > localMaxSpace {
						localBestSpaceHaltingMachine = newTm
					}

					localMaxNbSteps = MaxI(localMaxNbSteps, steps_count)
					localMaxSpace = MaxI(localMaxSpace, space_count)
					localNbHalt += 1

					if slow_down == 0 {
						wg.Add(1)

						go func() {
							Enumerate(nbStates, newTm, after_state, after_read, steps_count, space_count,
								SlowDownInit, simulation_backend)
							wg.Done()
						}()
					} else {
						Enumerate(nbStates, newTm, after_state, after_read, steps_count, space_count,
							slow_down-1, simulation_backend)
					}
					break

				case NO_HALT:
					localNbNoHalt += 1
					break

				case UNDECIDED_TIME:
					localNbUndecidedTime += 1
					UndecidedTimeLog.Write(newTm[:])
					break

				case UNDECIDED_SPACE:
					localNbUndecidedSpace += 1
					UndecidedSpaceLog.Write(newTm[:])
					break
				}
			}

		}

	}
	wg.Wait()

	mutexMetrics.Lock()
	NbMachineSeen += localNbMachineSeen
	NbMachinePruned += localNbMachinePruned
	NbHaltingMachines += localNbHalt
	NbNonHaltingMachines += localNbNoHalt
	NbUndecidedTime += localNbUndecidedTime
	NbUndecidedSpace += localNbUndecidedSpace

	if localMaxNbSteps >= MaxNbSteps {
		BBRecordLog.Write([]byte(fmt.Sprintf("*TIME %d SPACE %d\n%s\n",
			localMaxNbSteps, localMaxSpace,
			localBestTimeHaltingMachine.ToAsciiTable(nbStates))))
	} else if localMaxSpace >= MaxSpace {
		BBRecordLog.Write([]byte(fmt.Sprintf("TIME %d *SPACE %d\n%s\n",
			localMaxNbSteps, localMaxSpace,
			localBestSpaceHaltingMachine.ToAsciiTable(nbStates))))
	}

	MaxNbSteps = MaxI(localMaxNbSteps, MaxNbSteps)
	MaxSpace = MaxI(localMaxSpace, MaxSpace)
	MaxNbGoRoutines = MaxI(MaxNbGoRoutines, runtime.NumGoroutine())

	if Verbose && (!notFirstLog || time.Since(lastLogTime) >= time.Duration(LogFreq)) {
		notFirstLog = true
		lastLogTime = time.Now()
		fmt.Printf("run time: %s\ntotal: %d\npruned: %d (%.2f)\nhalt: %d (%.2f)\nnon halt: %d (%.2f)\nundecided time: %d (%.2f)\n"+
			"undecided space: %d (%.2f)\nbb est.: %d\nbb space est.: %d\nrun/sec: %f\nmax go routines: %d\n\n",
			time.Since(TimeStart), NbMachineSeen,
			NbMachinePruned, float64(NbMachinePruned)/float64(NbMachineSeen),
			NbHaltingMachines, float64(NbHaltingMachines)/float64(NbMachineSeen),
			NbNonHaltingMachines, float64(NbNonHaltingMachines)/float64(NbMachineSeen),
			NbUndecidedTime, float64(NbUndecidedTime)/float64(NbMachineSeen),
			NbUndecidedSpace, float64(NbUndecidedSpace)/float64(NbMachineSeen),
			MaxNbSteps, MaxSpace, float64(NbMachineSeen)/time.Since(TimeStart).Seconds(),
			MaxNbGoRoutines)

	}

	mutexMetrics.Unlock()

}
