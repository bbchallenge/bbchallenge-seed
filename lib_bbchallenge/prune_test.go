// Here we test our TM's filters
package bbchallenge

import "testing"

func TestPruneEquivalentStates(t *testing.T) {

	tm1 := TM{
		1, R, 3, 1, L, 3, // These states are equivalent
		1, R, 3, 1, L, 3, //
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}

	tm2 := TM{
		1, R, 2, 1, L, 1, // These states are equivalent
		1, R, 2, 1, L, 2, //
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}

	tm3 := TM{
		1, R, 3, 1, L, 1, // These states are equivalent
		1, R, 3, 1, L, 2, //
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}

	tm4 := TM{
		1, R, 2, 1, L, 3, // These states are NOT equivalent
		1, R, 3, 1, R, 2, //
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0}

	if !pruneEquivalentStates(5, tm1, 1) {
		t.Fail()
	}

	if !pruneEquivalentStates(5, tm2, 1) {
		t.Fail()
	}

	if !pruneEquivalentStates(5, tm3, 1) {
		t.Fail()
	}

	if pruneEquivalentStates(5, tm4, 1) {
		t.Fail()
	}
}
