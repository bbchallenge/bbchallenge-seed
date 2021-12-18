// Here we test our TMs simulation algorithms and various other things
package bbchallenge

import (
	"testing"
	"time"
)

func getBB5Winner() TM {
	// +---+-----+-----+
	// | - |  0  |  1  |
	// +---+-----+-----+
	// | A | 1RB | 1LC |
	// | B | 1RC | 1RB |
	// | C | 1RD | 0LE |
	// | D | 1LA | 1LD |
	// | E | 1RH | 0LA |
	// +---+-----+-----+

	return TM{
		1, R, 2, 1, L, 3,
		1, R, 3, 1, R, 2,
		1, R, 4, 0, L, 5,
		1, L, 1, 1, L, 4,
		1, R, 6, 0, L, 1}

}

func TestTabulateTM(t *testing.T) {

	bb5_winner := getBB5Winner()
	t.Log("\n" + bb5_winner.ToAsciiTable(5))

	notFullyDefinedTM := TM{
		1, R, 2, 1, L, 3,
		1, R, 3, 1, R, 2,
		1, R, 4, 0, L, 5,
		1, L, 1, 1, L, 4,
		0, 0, 0, 0, 0, 0}

	t.Log("\n" + notFullyDefinedTM.ToAsciiTable(5))
}

func TestBackendGo(t *testing.T) {
	start := time.Now()
	bb5_winner := getBB5Winner()

	halt_status, end_state, read, steps_count, space_count := simulate(bb5_winner, BB5, BB5_SPACE)

	if halt_status != HALT || end_state != H || read != 0 || steps_count != BB5 || space_count != BB5_SPACE {
		t.Error(halt_status, end_state, read, steps_count, space_count)
	}

	t.Log(time.Since(start))
}

func TestBackendC(t *testing.T) {
	start := time.Now()
	bb5_winner := getBB5Winner()
	halt_status, end_state, read, steps_count, space_count := simulate_C_wrapper(bb5_winner, BB5, BB5_SPACE)

	if halt_status != HALT || end_state != H || read != 0 || steps_count != BB5 || space_count != BB5_SPACE {
		t.Error(halt_status, end_state, read, steps_count, space_count)
	}
	t.Log(time.Since(start))
}
