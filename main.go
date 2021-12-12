package main

import (
	"flag"
	"fmt"
	"time"

	bbc "github.com/bbchallenge/bbchallenge/lib_bbchallenge"
)

const DEBUG = 1

var start time.Time

func main() {
	arg_nbStates := flag.Int("n", 4, "# of states")
	arg_backend := flag.Int("b", 0, "simulation backend (0 for go, 1 for C)")
	arg_verb := flag.Int("v", 0, "verbosity level (0 for no logs, 1 for final report and 2 for intermediate reports)")

	flag.Parse()

	bbc.TimeStart = time.Now()

	// Initial transition is 1RB (w.l.o.g)
	kick_start := bbc.TM{
		1, bbc.R, 2, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
	}

	nbStates := byte(*arg_nbStates)
	simulationBackend := bbc.SimulationBackend(*arg_backend)
	bbc.VERBOSITY = *arg_verb

	bbc.Search(nbStates, kick_start, 2, 0, 1, 1, 2, 2, simulationBackend)

	if bbc.VERBOSITY >= 1 {
		fmt.Println("Report")
		fmt.Println("======")
		fmt.Printf("Number of %d-state machines seen: %d\n", nbStates, bbc.NbMachineSeen)
		fmt.Printf("BB%d estimate: %d\n", nbStates, bbc.MaxNbSteps)
		fmt.Printf("BB%d_SPACE estimate: %d\n", nbStates, bbc.MaxSpace)
		fmt.Println("Max # of simultaneous Go routines during search:", bbc.MaxNbGoRoutines)
	}
}
