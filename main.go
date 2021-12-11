package main

import (
	"flag"
	"fmt"
	"time"
)

const DEBUG = 1

var VERBOSITY int

var nbStates byte

var start time.Time

func main() {
	arg_nbStates := flag.Int("n", 4, "# of states")
	arg_backend := flag.Int("b", 0, "simulation backend")
	arg_verb := flag.Int("v", 0, "verbosity level")
	flag.Parse()

	start = time.Now()

	// Initial transition is 1RB (w.l.o.g)
	kick_start := TM{
		1, R, 2, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}

	nbStates = byte(*arg_nbStates)
	simulationBackend := SimulationBackend(*arg_backend)
	VERBOSITY = *arg_verb

	search(kick_start, 2, 0, 1, 1, 2, 2, simulationBackend)

	if VERBOSITY >= 1 {
		fmt.Println("Report")
		fmt.Println("======")
		fmt.Printf("Number of %d-state machines seen: %d\n", nbStates, nbMachineSeen)
		fmt.Printf("BB%d estimate: %d\n", nbStates, maxNbSteps)
		fmt.Printf("BB%d_SPACE estimate: %d\n", nbStates, maxSpace)
		fmt.Println("Max # of simultaneous Go routines during search:", maxNbGoRoutines)
	}
}
