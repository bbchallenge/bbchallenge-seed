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

var DunnoTimeLog io.Writer
var DunnoSpaceLog io.Writer

var VERBOSE bool
var LOG_FREQ int64 = 30000000000
var TimeStart time.Time
var lastLogTime time.Time
var notFirstLog bool

var LIMIT_TIME int = BB5
var LIMIT_SPACE int = BB5_SPACE

var mutexMetrics sync.Mutex
var NbMachineSeen int
var NbHaltingMachines int
var NbNonHaltingMachines int
var NbDunnoTime int
var NbDunnoSpace int
var MaxNbSteps int
var MaxSpace int
var MaxNbGoRoutines int

// Invariant: tm's transition (state, read) is not defined
func Search(nbStates byte, tm TM, state byte, read byte,
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
	var localNbHalt int
	var localNbNoHalt int
	var localNbDunnoTime int
	var localNbDunnoSpace int
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
					haltStatus, after_state, after_read, steps_count, space_count = simulate(newTm, LIMIT_TIME, LIMIT_SPACE)
					break
				case SIMULATION_C:
					haltStatus, after_state, after_read, steps_count, space_count = simulate_C_wrapper(newTm, LIMIT_TIME, LIMIT_SPACE)
					break
				}

				switch haltStatus {
				case HALT:
					localMaxNbSteps = MaxI(localMaxNbSteps, steps_count)
					localMaxSpace = MaxI(localMaxSpace, space_count)
					localNbHalt += 1

					if slow_down == 0 {
						wg.Add(1)

						go func() {
							Search(nbStates, newTm, after_state, after_read, steps_count, space_count,
								slow_down_init, slow_down_init, simulation_backend)
							wg.Done()
						}()
					} else {
						Search(nbStates, newTm, after_state, after_read, steps_count, space_count,
							slow_down_init, slow_down-1, simulation_backend)
					}
					break

				case NO_HALT:
					localNbNoHalt += 1
					break

				case DUNNO_TIME:
					localNbDunnoTime += 1
					DunnoTimeLog.Write(tm[:])
					break

				case DUNNO_SPACE:
					localNbDunnoSpace += 1
					DunnoSpaceLog.Write(tm[:])
					break
				}
			}

		}

	}
	wg.Wait()

	mutexMetrics.Lock()
	NbMachineSeen += localNbMachineSeen

	NbHaltingMachines += localNbHalt
	NbNonHaltingMachines += localNbNoHalt
	NbDunnoTime += localNbDunnoTime
	NbDunnoSpace += localNbDunnoSpace

	MaxNbSteps = MaxI(localMaxNbSteps, MaxNbSteps)
	MaxSpace = MaxI(localMaxSpace, MaxSpace)
	MaxNbGoRoutines = MaxI(MaxNbGoRoutines, runtime.NumGoroutine())

	if VERBOSE && (!notFirstLog || time.Since(lastLogTime) >= time.Duration(LOG_FREQ)) {
		notFirstLog = true
		lastLogTime = time.Now()
		fmt.Printf("run time: %s\ntotal: %d\nhalt: %d (%.2f)\nnon halt: %d (%.2f)\ndunno time: %d (%.2f)\n"+
			"dunno space: %d (%.2f)\nbb est.: %d\nbb space est.: %d\nrun/sec: %f\nmax go routines: %d\n\n",
			time.Since(TimeStart), NbMachineSeen,
			NbHaltingMachines, float64(NbHaltingMachines)/float64(NbMachineSeen),
			NbNonHaltingMachines, float64(NbNonHaltingMachines)/float64(NbMachineSeen),
			NbDunnoTime, float64(NbDunnoTime)/float64(NbMachineSeen),
			NbDunnoSpace, float64(NbDunnoSpace)/float64(NbMachineSeen),
			MaxNbSteps, MaxSpace, float64(NbMachineSeen)/time.Since(TimeStart).Seconds(),
			MaxNbGoRoutines)

	}

	mutexMetrics.Unlock()

}
