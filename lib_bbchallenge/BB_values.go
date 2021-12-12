package bbchallenge

// Source: https://webusers.imj-prg.fr/~pascal.michel/bbc.html

// BB2 Winner:
// +---+-----+-----+
// | - |  0  |  1  |
// +---+-----+-----+
// | A | 1RB | 1LB |
// | B | 1LA | 1RH |
// +---+-----+-----+

const BB2 = 6
const BB2_SPACE = 4

// BB3 Winner:
// +---+-----+-----+
// | - |  0  |  1  |
// +---+-----+-----+
// | A | 1RB | 1RH |
// | B | 1LB | 0RC |
// | C | 1LC | 1LA |
// +---+-----+-----+

const BB3 = 21
const BB3_SPACE = 7

// BB4 Winner:
// +---+-----+-----+
// | - |  0  |  1  |
// +---+-----+-----+
// | A | 1RB | 1LB |
// | B | 1LA | 0LC |
// | C | 1RH | 1LD |
// | D | 1RD | 0RA |
// +---+-----+-----+

const BB4 = 107
const BB4_SPACE = 16

// BB5 Current Winner:
// +---+-----+-----+
// | - |  0  |  1  |
// +---+-----+-----+
// | A | 1RB | 1LC |
// | B | 1RC | 1RB |
// | C | 1RD | 0LE |
// | D | 1LA | 1LD |
// | E | 1RH | 0LA |
// +---+-----+-----+

// conjecture
const BB5 = 47176870
const BB5_SPACE = 12289
