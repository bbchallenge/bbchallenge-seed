package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	bbc "github.com/bbchallenge/bbchallenge/lib_bbchallenge"
	bbchallenge "github.com/bbchallenge/bbchallenge/lib_bbchallenge"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const DEBUG = 1

type BBChallengeFormatter struct {
}

func (f *BBChallengeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.

	return []byte(entry.Message + "\n"), nil
}

func getRunName() string {
	id, _ := uuid.NewV4()

	split := strings.Split(id.String(), "-")
	return "run-" + split[len(split)-1]
}

var dunnoTimeFile *os.File
var dunnoSpaceFile *os.File
var bbRecordFile *os.File

func initLogger(runName string) {

	mainLogFileName := runName + ".txt"
	log.SetFormatter(new(BBChallengeFormatter))
	log.SetOutput(bbchallenge.InitAppendFile(mainLogFileName, "output/"))

	dunnoTimeLogFileName := runName + "_dunno_time" // binary file
	bbc.DunnoTimeLog = bbchallenge.InitAppendFile(dunnoTimeLogFileName, "output/")

	dunnoSpaceLogFileName := runName + "_dunno_space" // binary file
	bbc.DunnoSpaceLog = bbchallenge.InitAppendFile(dunnoSpaceLogFileName, "output/")

	bbRecordLogFileName := runName + "_bb_records.txt"
	bbc.BBRecordLog = bbchallenge.InitAppendFile(bbRecordLogFileName, "output/")
}

func main() {
	runName := getRunName()
	initLogger(runName)

	arg_nbStates := flag.Int("n", 4, "# of states")
	arg_backend := flag.Int("b", 0, "simulation backend (0 for go, 1 for C)")
	arg_verb := flag.Bool("v", false, "displays infos about the current run on stdout")
	arg_verb_freq := flag.Int("vf", 30, "seconds between each stdout log in verbose mode")

	arg_limit_time := flag.Int("tlim", bbc.BB5, "time limit after which running machines are killed and marked as 'DUNNO_TIME' (known values of Busy Beaver are also used for early termination)")
	arg_limit_space := flag.Int("slim", bbc.BB5_SPACE, "space limit after which machines are killed and marked as 'DUNNO_SPACE' (known values of Busy Beaver space are also used for early termination)")

	arg_task_divisor := flag.Int("divtask", 1, "divides the size of the job by 1, 2, 4 or 8")

	arg_disable_filtering := flag.Bool("nf", false, "disable extra pruning of redundant machines from the enumeration")

	if !(*arg_task_divisor == 1 || *arg_task_divisor == 2 || *arg_task_divisor == 4 || *arg_task_divisor == 8) {

		fmt.Println("Task divisor must be either 1, 2, 4 or 8. Default is 1.")
		os.Exit(-1)
	}

	arg_task_divisor_me := flag.Int("mytask", 0, "select which task bucket this run will do")

	if *arg_task_divisor_me < 0 || *arg_task_divisor_me >= *arg_task_divisor {
		fmt.Println("Your task id must be either >= 0 and < the task divisor which is", *arg_task_divisor)
		os.Exit(-1)
	}

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

	log.Info(runName)
	log.Info(time.Now().Format(time.RFC1123))
	log.Info("Nb states: ", nbStates)

	bbc.Verbose = *arg_verb
	bbc.LogFreq = int64(*arg_verb_freq) * 1e9
	bbc.SimulationLimitTime = *arg_limit_time
	bbc.SimulationLimitSpace = *arg_limit_space
	bbc.SlowDownInit = 2
	bbc.ActivateFiltering = !*arg_disable_filtering

	bbc.TaskDivisor = *arg_task_divisor
	bbc.TaskDivisorMe = *arg_task_divisor_me

	log.Info("Task divisor: ", bbc.TaskDivisor)
	log.Info("My task: ", bbc.TaskDivisorMe)

	log.Info("Limit time: ", bbc.SimulationLimitTime)
	log.Info("Limit space: ", bbc.SimulationLimitSpace)

	if simulationBackend == bbc.SIMULATION_GO {
		log.Info("Simulation backend: GO")
	} else {
		log.Info("Simulation backend: C")
	}

	bbc.Enumerate(nbStates, kick_start, 2, 0, 1, 1, bbc.SlowDownInit, simulationBackend)

	log.Infoln("\nReport")
	log.Infoln("======")

	log.Info("Run time: ", time.Since(bbc.TimeStart), "\n")
	log.Info(fmt.Sprintf("Number of %d-state machines seen: %d", nbStates, bbc.NbMachineSeen))
	log.Info(fmt.Sprintf("Number of %d-state machines pruned: %d (%.2f)", nbStates, bbc.NbMachinePruned, float64(bbc.NbMachinePruned)/float64(bbc.NbMachineSeen)))
	log.Info(fmt.Sprintf("Number of halting machines: %d (%.2f)", bbc.NbHaltingMachines, float64(bbc.NbHaltingMachines)/float64(bbc.NbMachineSeen)))
	log.Info(fmt.Sprintf("Number of non-halting machines: %d (%.2f)", bbc.NbNonHaltingMachines, float64(bbc.NbNonHaltingMachines)/float64(bbc.NbMachineSeen)))
	log.Info(fmt.Sprintf("Number of dunno-time machines: %d (%.2f)", bbc.NbDunnoTime, float64(bbc.NbDunnoTime)/float64(bbc.NbMachineSeen)))
	log.Info(fmt.Sprintf("Number of dunno-space machines: %d (%.2f)\n", bbc.NbDunnoSpace, float64(bbc.NbDunnoSpace)/float64(bbc.NbMachineSeen)))

	log.Info(fmt.Sprintf("BB%d estimate: %d", nbStates, bbc.MaxNbSteps))
	log.Info(fmt.Sprintf("BB%d_SPACE estimate: %d\n", nbStates, bbc.MaxSpace))

	log.Info("Max # of simultaneous Go routines during search: ", bbc.MaxNbGoRoutines)
	log.StandardLogger().Writer().Close()

	dunnoTimeFile.Close()
	dunnoSpaceFile.Close()
	bbRecordFile.Close()
}
