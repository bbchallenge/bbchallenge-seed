package main

// We currently work with machines that have at most MAX_STATES states
const MAX_STATES = 5

// Name of halting state
const H = 6

const BASE = 2
const SEP = 3

const MAX_MEMORY = 40000

type HaltStatus byte

const (
	HALT HaltStatus = iota
	NO_HALT
	DUNNO_TIME
	DUNNO_SPACE
)

func MaxI(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinI(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// We are considering <= 5-state 2-symbol TMs
// The TM:
//
// +---+-----------+-----+
// | - |     0     |  1  |
// +---+-----------+-----+
// | A | 1RB       | 1RH |
// | B | 1LB       | 0RC |
// | C | 1LC       | 1LA |
// | D | undefined | 1RA |
// +---+-----------+-----+
//
// Is encoded by the array:
// 1, 0, 2, 1, 1, 6, 1, 1, 2, 0, 0, 3, 1, 1, 3  1, 1, 1, 0, 0, 0, 1, 0, 1
// 1, R, B, 1, R, H, 1, L, B, 0, R, C, 1, L, C, 1, L, A, -, -, -, 1, R, A

type TM [2 * MAX_STATES * 3]byte

// Simulates the input TM from blank input
// and state 1.
// Returns undetermined, state, read with:
// - halting status (HaltStatus)
// - state (byte): State number of undetermined transition if reached
// - read (byte): Read symbol of undetermined transition if reached
// - steps count
// - space count
func simulate(tm TM) (HaltStatus, byte, byte, int, int) {
	var tape [MAX_MEMORY]byte

	max_pos := 0
	min_pos := MAX_MEMORY - 1
	curr_head := MAX_MEMORY / 2

	var curr_state byte = 1

	steps_count := 0

	var state_seen [MAX_STATES]bool
	var nbStateSeen byte

	for curr_state != H {

		if !state_seen[curr_state-1] {
			nbStateSeen += 1
		}
		state_seen[curr_state-1] = true

		// Using knowledge about BB time and space
		if nbStateSeen <= 4 && steps_count > BB4 {
			return NO_HALT, 0, 0,
				steps_count, max_pos - min_pos + 1
		}

		if nbStateSeen <= 4 && max_pos-min_pos+1 > BB4_SPACE {
			return NO_HALT, 0, 0,
				steps_count, max_pos - min_pos + 1
		}

		if nbStateSeen == 5 && steps_count > BB5 {
			return DUNNO_TIME, 0, 0, steps_count, max_pos - min_pos + 1
		}

		if nbStateSeen == 5 && max_pos-min_pos+1 > BB5_SPACE {
			return DUNNO_SPACE, 0, 0, steps_count, max_pos - min_pos + 1
		}

		min_pos = MinI(min_pos, curr_head)
		max_pos = MaxI(max_pos, curr_head)

		read := tape[curr_head]

		tm_transition := 6*(curr_state-1) + 3*read
		write := tm[tm_transition]
		move := tm[tm_transition+1]
		next_state := tm[tm_transition+2]

		// undefined transition
		if next_state == 0 {
			return HALT, curr_state, read,
				steps_count + 1, max_pos - min_pos + 1
		}

		tape[curr_head] = write

		if move == R {
			curr_head += 1
			if curr_head == MAX_MEMORY {
				return DUNNO_SPACE, 0, 0,
					steps_count, max_pos - min_pos + 1
			}

		} else {
			curr_head -= 1
			if curr_head == -1 {
				return DUNNO_SPACE, 0, 0,
					steps_count, max_pos - min_pos + 1
			}
		}

		steps_count += 1
		curr_state = next_state
	}

	return HALT, H, 0,
		steps_count, max_pos - min_pos + 1
}
