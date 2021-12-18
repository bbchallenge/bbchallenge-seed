// Here we simulate TMs in C
#include "simulate.h"

#define MAX_MEMORY 40000
#define MAX_STATES 5

#define H 6
#define R 0
#define L 1

#define BB4 107
#define BB4_SPACE 16

// conjecture
#define BB5 47176870
#define BB5_SPACE 12289

#define RETURN(HALT_STATUS, STATE, READ, STEPS_COUNT, SPACE_COUNT) \
  {                                                                \
    *ret_state = STATE;                                            \
    *ret_read = READ;                                              \
    *ret_steps_count = STEPS_COUNT;                                \
    *ret_space_count = SPACE_COUNT;                                \
    return HALT_STATUS;                                            \
  }

typedef unsigned char BYTE;

const BYTE HALT = 0;
const BYTE NO_HALT = 1;
const BYTE DUNNO_TIME = 2;
const BYTE DUNNO_SPACE = 3;

BYTE simulate(BYTE* tm,
              int limit_time,
              int limit_space,
              BYTE* ret_state,
              BYTE* ret_read,
              int* ret_steps_count,
              int* ret_space_count) {
  BYTE tape[MAX_MEMORY] = {0};

  int max_pos = 0;
  int min_pos = MAX_MEMORY - 1;
  int curr_head = MAX_MEMORY / 2;

  BYTE curr_state = 1;
  int steps_count = 0;

  BYTE state_seen[MAX_STATES] = {0};
  BYTE nbStateSeen = 0;

  BYTE read;

  while (curr_state != H) {
    if (!state_seen[curr_state - 1]) {
      nbStateSeen += 1;
    }
    state_seen[curr_state - 1] = 1;

    if (nbStateSeen <= 4 && steps_count > BB4) {
      RETURN(NO_HALT, 0, 0, steps_count, max_pos - min_pos + 1)
    }

    if (nbStateSeen <= 4 && max_pos - min_pos + 1 > BB4_SPACE) {
      RETURN(NO_HALT, 0, 0, steps_count, max_pos - min_pos + 1)
    }

    if (nbStateSeen == 5 && steps_count > limit_time) {
      RETURN(DUNNO_TIME, 0, 0, steps_count, max_pos - min_pos + 1)
    }

    if (nbStateSeen == 5 && max_pos - min_pos + 1 > limit_space) {
      RETURN(DUNNO_SPACE, 0, 0, steps_count, max_pos - min_pos + 1)
    }

    if (curr_head < min_pos) {
      min_pos = curr_head;
    }

    if (curr_head > max_pos) {
      max_pos = curr_head;
    }

    read = tape[curr_head];

    BYTE tm_transition = 6 * (curr_state - 1) + 3 * read;
    BYTE write = tm[tm_transition];
    BYTE move = tm[tm_transition + 1];
    BYTE next_state = tm[tm_transition + 2];

    // undefined transition
    if (next_state == 0) {
      RETURN(HALT, curr_state, read, steps_count + 1, max_pos - min_pos + 1)
    }

    tape[curr_head] = write;

    if (move == R) {
      curr_head += 1;
      if (curr_head == MAX_MEMORY) {
        RETURN(DUNNO_SPACE, 0, 0, steps_count, max_pos - min_pos + 1)
      }
    } else {
      curr_head -= 1;
      if (curr_head == -1) {
        RETURN(DUNNO_SPACE, 0, 0, steps_count, max_pos - min_pos + 1)
      }
    }

    steps_count += 1;
    curr_state = next_state;
  }

  RETURN(HALT, H, read, steps_count, max_pos - min_pos + 1);
}