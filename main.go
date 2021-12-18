package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	bbc "github.com/bbchallenge/bbchallenge/lib_bbchallenge"
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

func initAppendFile(logFileName string) *os.File {
	ioutil.WriteFile("output/"+logFileName, []byte(""), 0644)
	logFile, _ := os.OpenFile("output/"+logFileName, os.O_APPEND|os.O_WRONLY, 0644)
	return logFile
}

var dunnoTimeFile *os.File
var dunnoSpaceFile *os.File
var bbRecordFile *os.File

func initLogger(runName string) {

	mainLogFileName := runName + ".txt"
	log.SetFormatter(new(BBChallengeFormatter))
	log.SetOutput(initAppendFile(mainLogFileName))

	dunnoTimeLogFileName := runName + "_dunno_time" // binary file
	bbc.DunnoTimeLog = initAppendFile(dunnoTimeLogFileName)

	dunnoSpaceLogFileName := runName + "_dunno_space" // binary file
	bbc.DunnoSpaceLog = initAppendFile(dunnoSpaceLogFileName)

	bbRecordLogFileName := runName + "_bb_records.txt"
	bbc.BBRecordLog = initAppendFile(bbRecordLogFileName)
}

func main() {
	runName := getRunName()
	initLogger(runName)

	arg_nbStates := flag.Int("n", 4, "# of states")
	arg_backend := flag.Int("b", 0, "simulation backend (0 for go, 1 for C)")
	arg_verb := flag.Bool("v", false, "displays infos about the run every 5m on stdout")
	arg_verb_freq := flag.Int("vf", 30, "seconds between each stdout log in verbose mode")

	arg_limit_time := flag.Int("lt", bbc.BB5, "Time limit after which running machines are killed and marked as 'DUNNO_TIME' (known values of Busy Beaver are also used for early termination)")
	arg_limit_space := flag.Int("ls", bbc.BB5_SPACE, "Space limit after which machines are killed and marked as 'DUNNO_SPACE' (known values of Busy Beaver space are also used for early termination)")

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

	log.Info("Limit time: ", bbc.SimulationLimitTime)
	log.Info("Limit space: ", bbc.SimulationLimitSpace)

	if simulationBackend == bbc.SIMULATION_GO {
		log.Info("Simulation backend: GO")
	} else {
		log.Info("Simulation backend: C")
	}

	bbc.Search(nbStates, kick_start, 2, 0, 1, 1, bbc.SlowDownInit, simulationBackend)

	log.Infoln("\nReport")
	log.Infoln("======")

	log.Info("Run time: ", time.Since(bbc.TimeStart), "\n")
	log.Info(fmt.Sprintf("Number of %d-state machines seen: %d", nbStates, bbc.NbMachineSeen))
	log.Info(fmt.Sprintf("Number of halting machines: %d", bbc.NbHaltingMachines))
	log.Info(fmt.Sprintf("Number of non-halting machines: %d", bbc.NbNonHaltingMachines))
	log.Info(fmt.Sprintf("Number of dunno-time machines: %d", bbc.NbDunnoTime))
	log.Info(fmt.Sprintf("Number of dunno-spqce machines: %d\n", bbc.NbDunnoSpace))

	log.Info(fmt.Sprintf("BB%d estimate: %d", nbStates, bbc.MaxNbSteps))
	log.Info(fmt.Sprintf("BB%d_SPACE estimate: %d\n", nbStates, bbc.MaxSpace))

	log.Info("Max # of simultaneous Go routines during search: ", bbc.MaxNbGoRoutines)
	log.StandardLogger().Writer().Close()

	dunnoTimeFile.Close()
	dunnoSpaceFile.Close()
	bbRecordFile.Close()
}
